package repository

import (
	"context"
	"github.com/Asylann/gRPC_Demo/server/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type sqlCartStore interface {
	GetCartById(ctx context.Context, id int) (models.Cart, error)
	CreateCart(ctx context.Context, userId int) error
	AddToCart(ctx context.Context, cartId int, product models.Product) error
}

type CartStore struct {
	getById    *sqlx.Stmt
	createCart *sqlx.Stmt
	AddTo      *sqlx.Stmt
	db         *sqlx.DB
}

func NewCartStore(db *sqlx.DB) (CartStore, error) {
	cartstore := CartStore{db: db}

	var err error
	cartstore.getById, err = db.PreparexContext(context.Background(),
		`SELECT * FROM carts WHERE ID=$1`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	cartstore.createCart, err = db.PreparexContext(context.Background(),
		`INSERT INTO carts (user_id) VALUES($1)`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	cartstore.AddTo, err = db.PreparexContext(context.Background(),
		`INSERT INTO cart_items (cart_id,product_id) VALUES($1,$2)`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	return cartstore, nil
}

func (c *CartStore) GetCartById(id int) (models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cart models.Cart
	if err := c.getById.GetContext(ctx, &cart, id); err != nil {
		log.Println(err.Error())
		return models.Cart{}, err
	}
	return cart, nil
}

func (c *CartStore) CreateCart(userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := c.createCart.ExecContext(ctx, userId); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (c *CartStore) AddToCart(cardId int, product models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := c.AddTo.ExecContext(ctx, cardId, product.ID); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
