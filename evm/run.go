package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"inscription-mint-bot/evm/config"
	"math/big"
	"strings"
	"sync"
	"time"
)

var configFile = flag.String("f", "etc/deploy.yaml", "the config file")

func main() {
	//读取配置文件
	flag.Parse()

	var c config.Config
	var wg sync.WaitGroup

	conf.MustLoad(*configFile, &c)
	logx.Infof("===============read config success===============")

	priKeys := c.PriKeys
	if len(priKeys) < 1 {
		logx.Errorf(" no prikey ")
		return
	}
	defer ants.Release()

	for _, priKey := range priKeys {
		wg.Add(1)
		priKeyLocal := priKey
		ants.Submit(func() {

			fromAddress, err := getAddress(priKeyLocal)
			if err != nil {
				return
			}
			//获取发送数据
			inputData := strings.TrimSpace(c.MintConf.InputData)
			if inputData == "" {
				return
			}
			//获取发送地址,默认不填,不填给自己地址发送,填了给指定地址发送
			toAddress := fromAddress
			if c.MintConf.ToAddr != "" {
				tmpToddress := common.HexToAddress(c.MintConf.ToAddr)
				toAddress = &tmpToddress
			}
			//获取rpc客户端
			client, err := ethclient.Dial(c.EthRpcConf.Url)
			if err != nil {
				logx.Errorf(" ethclient.Dial:%s", err.Error())
				return
			}
			//获取chainID
			chainID, err := client.NetworkID(context.Background())
			if err != nil {
				logx.Errorf(" client.NetworkID:%s", err.Error())
			}
			//获取nonce
			nonce, err := client.PendingNonceAt(context.Background(), *fromAddress)
			if err != nil {
				logx.Errorf("crypto.PubkeyToAddress:%s", err.Error())
			}
			value := big.NewInt(0)
			gasLimit := uint64(210000)
			ratio := c.EthRpcConf.GasPriceRatio

			privateKey, err := crypto.HexToECDSA(priKeyLocal)

			//定时执行
			ticker := time.NewTicker(time.Duration(c.EthRpcConf.IntervalTime) * time.Millisecond)
			defer ticker.Stop()
			for range ticker.C {

				gasPrice, err := client.SuggestGasPrice(context.Background())

				if err != nil {
					logx.Errorf("client.SuggestGasPrice:%s", err.Error())
					continue
				}
				//获取gasPrice比例,如果大于1,修改gasPrice
				if ratio > 1 {

					ratioFloat := big.NewFloat(ratio)
					gasPriceInt64, _ := ratioFloat.Mul(new(big.Float).SetInt64(gasPrice.Int64()), ratioFloat).Int64()
					gasPrice = big.NewInt(gasPriceInt64)
				}
				//组装交易数据
				tx := types.NewTransaction(nonce, *toAddress, value, gasLimit, gasPrice, []byte(inputData))

				signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
				if err != nil {
					logx.Errorf(" client.NetworkID:%s", err.Error())
					continue
				}
				//发送交易
				err = client.SendTransaction(context.Background(), signedTx)
				if err != nil {
					logx.Errorf(" SendTransaction:%s", err.Error())
					continue
				}

				logx.Infof("fromAddress:%s,toAddress:%s,nonce:%d,gas:%v,txhash:%v", fromAddress, toAddress, nonce, gasPrice, signedTx.Hash())

				nonce++
			}

			wg.Done()
		})

	}
	wg.Wait()
}

// 根据私钥获取地址
func getAddress(priKey string) (*common.Address, error) {
	privateKey, err := crypto.HexToECDSA(priKey)
	if err != nil {
		logx.Errorf("crypto.HexToECDSA:%s", err.Error())
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logx.Errorf("ecdsa.PublicKey:%s", err.Error())
		return nil, err
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &address, nil
}
