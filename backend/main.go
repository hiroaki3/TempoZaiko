package TempoZaiko0x0

import (
	"net/http"
	"strconv"
)

// 初期化
func init() {
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.HandleFunc("/api/", restHandler)
}

// REST用のハンドラ
func restHandler(w http.ResponseWriter, r *http.Request) {

	p := NewParam(r.URL)
	if p.Kind != "expense" {
		respondErr(w, r, http.StatusBadRequest, "expense以外の種別には対応してません")
	}

	// PUT以外の場合はprice用にダミーの0を入れる
	if r.Method != "PUT" && !p.HasValue() {
		p.Value = "0"
	}

	// kindがexpenseの場合はvalueをintと定義したのでキャスト
	price, err := strconv.Atoi(p.Value)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ed := NewExpenseDatastore(r)
	if err := ed.build(p.Key, price); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err.Error())
	}

	switch r.Method {
	case "GET":
		handleGet(ed, w, r)
		return
	case "PUT":
		handlePut(ed, w, r)
		return
	case "DELETE":
		handleDelete(ed, w, r)
		return
	default:
		respondErr(w, r, http.StatusNotFound, "未対応のHTTPメソッドです")
	}
}

type SuccessResponse struct {
	Expense ExpenseEntry `json:"entity"`
	Message string       `json:"message"`
}

// GET用のhandler。エンティティの取得を行う
func handleGet(ed *ExpenseDatastore, w http.ResponseWriter, r *http.Request) {
	if err := ed.Get(); err != nil {
		respondErr(w, r, http.StatusBadRequest, err.Error())
		return
	}
	message := "「" + ed.Expense.Name + "」の金額は" + strconv.Itoa(ed.Expense.Price) + "円です。"
	respond(w, r, http.StatusOK, SuccessResponse{*ed.Expense, message})
}

// PUT用のhandler。エンティティの作成・更新を行う
func handlePut(ed *ExpenseDatastore, w http.ResponseWriter, r *http.Request) {
	if err := ed.Put(); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	message := "「" + ed.Expense.Name + "」の登録を行いました。"
	respond(w, r, http.StatusOK, SuccessResponse{*ed.Expense, message})
}

// DELETE用のhandler。エンティティの削除を行う
func handleDelete(ed *ExpenseDatastore, w http.ResponseWriter, r *http.Request) {
	if err := ed.Delete(); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	message := "「" + ed.Expense.Name + "」の削除を行いました。"
	respond(w, r, http.StatusOK, SuccessResponse{*ed.Expense, message})
}
