package repository

import (
	"context"
	"github.com/Asylann/gRPC_Demo/server/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type sqlCartStore interface {
	GetCartByUserId(ctx context.Context, id int) (models.Cart, error)
	CreateCart(ctx context.Context, userId int) error
	AddToCart(ctx context.Context, cartId int, product models.Product) error
}

type CartStore struct {
	getByUserId          *sqlx.Stmt
	createCart           *sqlx.Stmt
	AddTo                *sqlx.Stmt
	getProductOfCartById *sqlx.Stmt
	deleteItemFromCart   *sqlx.Stmt
	db                   *sqlx.DB
}

func NewCartStore() (CartStore, error) {
	cartstore := CartStore{db: db}

	var err error
	cartstore.getByUserId, err = db.PreparexContext(context.Background(),
		`SELECT * FROM carts WHERE user_id=$1`)
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

	cartstore.getProductOfCartById, err = db.PreparexContext(context.Background(),
		`SELECT product_id FROM cart_items WHERE cart_id = $1`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	cartstore.deleteItemFromCart, err = db.PreparexContext(context.Background(),
		`DELETE FROM cart_items WHERE cart_id=$1 and product_id=$2`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	return cartstore, nil
}

func (c *CartStore) GetCartByUserId(id int) (models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cart models.Cart
	if err := c.getByUserId.GetContext(ctx, &cart, id); err != nil {
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

func (c *CartStore) GetProductsOfCartById(cartId int) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var productsIds []int
	err := c.getProductOfCartById.SelectContext(ctx, &productsIds, cartId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return productsIds, nil
}

func (c *CartStore) DeleteItemFromCart(cartId, product_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := c.deleteItemFromCart.ExecContext(ctx, cartId, product_id); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
