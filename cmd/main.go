package main

import (
	"context"
	"fmt"
	"github.com/paintmeyellow/gotron/pkg/account"
	"github.com/paintmeyellow/gotron/pkg/client/grpc"
	"github.com/paintmeyellow/gotron/pkg/client/rest"
	"github.com/paintmeyellow/gotron/pkg/common/crypto"
	"log"
	"net/url"
	"time"
)

func main() {
	fmt.Println("----> Creating Account...")
	acc := createAccount()
	fmt.Println("<---- Done.")
	fmt.Println()

	grpcClient := grpc.NewClient("grpc.shasta.trongrid.io:50051")
	if err := grpcClient.Start(); err != nil {
		log.Fatal(err)
	}
	//Owner Public Key: TTbfNX2MbRJ1NG1FqQmKUy4E24sV9K2yzj
	//Owner Private Key: 9a5717a421b16aa5ddba886ba24849f09e78b47a0d1b8c67d1dd34c97bf31096
	fmt.Println("----> Transfering Transaction...")
	priv := "9a5717a421b16aa5ddba886ba24849f09e78b47a0d1b8c67d1dd34c97bf31096"
	transferTxn(grpcClient, priv, acc.Address, 1)
	fmt.Println("<---- Done.")
	fmt.Println()

	fmt.Println("----> Getting balance...")
	getBalance(grpcClient, "TTbfNX2MbRJ1NG1FqQmKUy4E24sV9K2yzj")
	fmt.Println("<---- Done.")
	fmt.Println()

	httpAddr, err := url.Parse("https://api.shasta.trongrid.io")
	if err != nil {
		log.Fatal(err)
	}
	restClient := rest.NewClient(*httpAddr)
	fmt.Println("----> Fetching Deposits List...")
	fetchDepositsList(restClient, 0, []string{
		//"TL8oHjkJTzqEPbKMRVuEVZNwzdgBNPhWd5",
		//"TBXUP7koqXbpbam2Dpt6rG1bxgAv3hRqAu",
		acc.Address,
	})
	fmt.Println("<---- Done.")
}

func createAccount() *account.Account {
	acc, err := account.Create()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tron Address:", acc.Address)
	fmt.Println("Private Key:", acc.PrivateKey)
	return acc
}

func getBalance(client *grpc.Client, addr string) {
	grpcAcc, err := client.GetAccount(context.Background(), addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Balance:", grpcAcc.GetBalance())
}

func transferTxn(client *grpc.Client, priv string, toAddr string, amount int64) {
	privKeyECDSA, err := crypto.HexToECDSA(priv)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := crypto.Base58ToAddress(toAddr)
	if err != nil {
		log.Fatal(err)
	}

	sendTxn := grpc.SendableTransaction{
		OwnerKey:  *privKeyECDSA,
		ToAddress: addr,
		Amount:    amount,
	}
	hash, err := client.Transfer(context.Background(), sendTxn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Transaction Hash:", hash.Hex())
	fmt.Println("<---- Done.")
	fmt.Println()

	fmt.Println("----> Fetching Transaction Info...")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()
	txn, err := client.GetTransactionFullInfoByID(ctx, hash, 5, time.Second)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash:", txn.Hash.Hex())
	fmt.Println("From:", txn.From.String())
	fmt.Println("To:", txn.To.String())
	fmt.Println("Amount:", txn.Amount)
	fmt.Println("NetFee:", txn.NetFee)
}

func fetchDepositsList(client *rest.Client, minTimestamp int64, addresses []string) {
	deposits, err := client.ListDeposits(context.Background(), minTimestamp, addresses)
	if err != nil {
		log.Fatal(err)
	}
	for _, deposit := range deposits {
		fmt.Println("Address:", deposit.Address)
		fmt.Println("Hash:", deposit.Hash)
		fmt.Println("Amount:", deposit.Amount)
		fmt.Println("Timestamp:", deposit.TransactionTimestamp)
		fmt.Println("Confirmations:", deposit.Confirmations)
	}
}
