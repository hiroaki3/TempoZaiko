package TempoZaiko0x0

import (
	"errors"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const ENTITYNAME = "Expense"

// 今回のサンプルに特化した、Datastoreに格納するエンティティ用の構造体
// jsonレスポンスとして利用するためタグ名を指定
type ExpenseEntry struct {
	Name  string `json:"string"`
	Price int    `json:"price"`
}

// ExpenseをDatastoreにCRUDするための構造体と手続き群。
// (名前が適当すぎてDatastoreの費用を管理してそうな名前に…)
// [使い方]
//  ed := NewExpenseDatasotre(r)
//  ed.build("項目名", 1080)
//  ed.Put()
type ExpenseDatastore struct {
	Ctx     context.Context
	Keys    []*datastore.Key
	Expense *ExpenseEntry
}

// コンストラクタ的な存在
func NewExpenseDatastore(r *http.Request) *ExpenseDatastore {
	ed := new(ExpenseDatastore)
	ed.Ctx = appengine.NewContext(r)
	return ed
}

// datastoreに渡す用の構造体の作成と、そのNameに一致するkeysの取得を行う
func (ed *ExpenseDatastore) build(name string, price int) error {
	ed.Expense = &ExpenseEntry{
		Name:  name,
		Price: price,
	}
	// Entityの中身の確認用に標準出力にログ出力。↓な感じのログが出ます。(フォーマットは%#vがお気に入り）
	// 2016/11/06 15:38:32 &hello_datastore.ExpenseEntiry{Name:"test", Price:55}
	log.Printf("%#v", ed.Expense)
	var err error
	ed.Keys, err = getKeysByName(ed)
	return err
}

// 構造体のポインタにエンティティが読み込まれる
// 今回の簡易サンプルでは問題ないが、フィールドが一致しない項目は上書きされないため注意
func (ed *ExpenseDatastore) Get() error {
	if len(ed.Keys) == 0 {
		return errors.New("取得対象の項目が存在しません")
	}
	return datastore.Get(ed.Ctx, ed.Keys[0], ed.Expense)
}

// このput処理ではputの複数同時リクエストがNameの一意性が失われてしまうが、今回のサンプルでは考慮しない
// 次回以降にgoのvarsLockの勉強と合わせて対応してみる予定
func (ed *ExpenseDatastore) Put() error {
	// keyがある場合は上書き、ない場合は新規作成
	if len(ed.Keys) == 0 {
		keys := make([]*datastore.Key, 1)
		// 今回のサンプルでは、新規作成時のStringIDはランダムでいいためIncompleteKeyを利用
		keys[0] = datastore.NewIncompleteKey(ed.Ctx, ENTITYNAME, nil)
		ed.Keys = keys
	}

	_, err := datastore.Put(ed.Ctx, ed.Keys[0], ed.Expense)
	return err
}

//　削除のみmultiに対応。putでnameの重複を発生させることで実験可能。
//　※トランザクションの勉強用
func (ed *ExpenseDatastore) Delete() error {
	if len(ed.Keys) == 0 {
		return errors.New("削除対象の項目が存在しません")
	}

	// 今回のサンプルではエンティティは全てROOTエンティティなため、複数削除にはクロスグループトランザクションを指定
	option := &datastore.TransactionOptions{XG: true}
	return datastore.RunInTransaction(ed.Ctx, func(c context.Context) error {
		return datastore.DeleteMulti(c, ed.Keys)
	}, option)
}

//　Nameと完全一致するEntityのKeyリストを取得
func getKeysByName(ed *ExpenseDatastore) ([]*datastore.Key, error) {
	q := datastore.NewQuery(ENTITYNAME).Filter("Name =", ed.Expense.Name)

	var expenses []ExpenseEntry
	return q.GetAll(ed.Ctx, &expenses)

}
