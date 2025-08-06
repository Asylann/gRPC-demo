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
	getByUserId                  *sqlx.Stmt
	createCart                   *sqlx.Stmt
	AddTo                        *sqlx.Stmt
	getProductOfCartById         *sqlx.Stmt
	deleteItemFromCart           *sqlx.Stmt
	getEtagVersionByUserId       *sqlx.Stmt
	changeEtagVersionByUserId    *sqlx.Stmt
	setDefaultEtagVersion        *sqlx.Stmt
	deleteCart                   *sqlx.Stmt
	deleteProductOfCarts         *sqlx.Stmt
	changeEtagVersionByProductId *sqlx.Stmt
	db                           *sqlx.DB
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
		`DELETE FROM cart_items 
       		   WHERE ctid IN (
       		       SELECT ctid
       		       FROM cart_items
       		       WHERE cart_id=$1 and product_id=$2
       		       LIMIT 1 
       		       );`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	cartstore.getEtagVersionByUserId, err = db.PreparexContext(context.Background(),
		`SELECT version FROM etag_versions WHERE userid = $1`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	cartstore.changeEtagVersionByUserId, err = db.PreparexContext(context.Background(),
		`UPDATE etag_versions SET version=$2 WHERE userid=$1`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	cartstore.setDefaultEtagVersion, err = db.PreparexContext(context.Background(),
		`INSERT INTO etag_versions (userid,version) VALUES($1,1)`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	cartstore.deleteCart, err = db.PreparexContext(context.Background(),
		`DELETE FROM carts WHERE user_id=$1 RETURNING id , user_id`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	cartstore.deleteProductOfCarts, err = db.PreparexContext(context.Background(),
		`DELETE FROM cart_items WHERE product_id=$1`)
	if err != nil {
		log.Println(err.Error())
		return CartStore{}, err
	}

	cartstore.changeEtagVersionByProductId, err = db.PreparexContext(context.Background(),
		`WITH affected_users AS ( 
    				SELECT DISTINCT c.user_id 
					FROM carts c
					JOIN cart_items ci ON c.id=ci.cart_id
					WHERE ci.product_id = $1
				)
   				UPDATE etag_versions
   				SET version=version+1
   				WHERE userid IN (SELECT user_id FROM affected_users)`,
	)
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

func (c *CartStore) GetEtagVersionByUserId(userId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var version int
	if err := c.getEtagVersionByUserId.GetContext(ctx, &version, userId); err != nil {
		_, err := c.setDefaultEtagVersion.ExecContext(ctx, userId)
		log.Println("Default was set of userId", userId)
		if err != nil {
			log.Println(err.Error())
			return 0, err
		}
	}
	if err := c.getEtagVersionByUserId.GetContext(ctx, &version, userId); err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return version, nil
}

func (c *CartStore) ChangeEtagVersionByUserId(userId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	version, err := c.GetEtagVersionByUserId(userId)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	version++

	if _, err := c.changeEtagVersionByUserId.ExecContext(ctx, userId, version); err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return version, nil
}

func (c *CartStore) DeleteCart(user_id int) (models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var cart models.Cart
	if err := c.deleteCart.QueryRowContext(ctx, user_id).Scan(&cart.Id, &cart.User_id); err != nil {
		log.Println(err.Error())
		return cart, err
	}

	return cart, nil
}

func (c *CartStore) DeleteProductOfCarts(productId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := c.deleteProductOfCarts.ExecContext(ctx, productId); err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Product by id = %v was deleted from all carts", productId)
	return nil
}

func (c *CartStore) ChangeEtagVersionByProductId(productId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := c.changeEtagVersionByProductId.ExecContext(ctx, productId); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
