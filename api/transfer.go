package api

// type transferRequest struct {
// 	FromAccountID int64  `json:"from_id" binding:"required, min=1"`
// 	ToAccountID   int64  `json:"to_id" binding:"required, min=1"`
// 	Amount        int64  `json:"amount" binding:"required,gt=0"`
// 	Currency      string `json:"currency" binding:"required,currency"`
// }

// func (server *Server) createEntry(ctx *gin.Context) {
// 	var req transferRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	arg := db.CreateEntryParams{
// 		AccountID: req.AccountID,
// 		Amount:    req.Amount,
// 	}

// 	entry, err := server.store.CreateEntry(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, entry)
// }

// func (server *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) {

// }
