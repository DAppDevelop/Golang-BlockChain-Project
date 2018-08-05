package main

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

func main() {

	//创建表
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//创建表
	//err = db.Update(func(tx *bolt.Tx) error {
	//	// 创建BlockBucket表
	//	//b, err := tx.CreateBucket([]byte("blockbucket"))
	//	//if err != nil {
	//	//	return fmt.Errorf("create bucket: %s", err)
	//	//}
	//	//
	//	//if b != nil {
	//	//	err := b.Put([]byte("ll"),[]byte("Send 100 BTC To 关系哥......"))
	//	//	if err != nil {
	//	//		log.Panic("数据存储失败......")
	//	//	}
	//	//	fmt.Println("存储成功")
	//	//}
	//	//
	//	//return nil
	//
	//	//获取BlockBucket表单
	//	b := tx.Bucket([]byte("blockbucket"))
	//
	//	// 往表里面存储数据
	//	if b != nil {
	//		err := b.Put([]byte("ll"),[]byte("Send 1000 BTC To 冠希哥......"))
	//		if err != nil {
	//			log.Panic("数据存储失败......")
	//		}
	//		fmt.Println("成功")
	//	}
	//
	//	//返回nil，以便数据库处理相应操作
	//	return nil
	//
	//})

	//查看数据
	err = db.View(func(tx *bolt.Tx) error {

		// 获取BlockBucket表对象
		b := tx.Bucket([]byte("blockbucket"))

		// 往表里面存储数据
		if b != nil {
			data := b.Get([]byte("l"))
			fmt.Printf("%s\n",data)
			data = b.Get([]byte("ll"))
			fmt.Printf("%s\n",data)
		}

		// 返回nil，以便数据库处理相应操作
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}


}
