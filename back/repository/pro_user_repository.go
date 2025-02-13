package repository

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gulmarket/model"
    "github.com/jmoiron/sqlx"
    "github.com/lib/pq"
    "log"
)

type pGProUserRepository struct {
    DB *sqlx.DB
}

func NewProUserRepository(db *sqlx.DB) model.ProUserRepository {
    return &pGProUserRepository{
        DB: db,
    }
}

func (r *pGProUserRepository) Create(ctx context.Context, p *model.ProUser) error {
    query := `INSERT INTO pro_users (email, password, phone, role)
              VALUES ($1, $2, $3, $4) RETURNING *`

    if err := r.DB.GetContext(ctx, p, query, p.Email, p.Password, p.Phone, p.Role); err != nil {
        if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
            log.Printf("Could not create a Pro User with email: %v. Reason: %v\n", p.Email, err.Code.Name())
            return err
        }

        log.Printf("Could not create a Pro User with email: %v. Reason: %v\n", p.Email, err)
        return err
    }
    return nil
}

func (r *pGProUserRepository) UpdateProUserRow(ctx context.Context, value interface{}, row string, id int) error {
    allowedColumns := map[string]bool{
        "email":    true,
        "password": true,
        "name":     true,
        "phone":    true,
        "city":     true,
    }

    if _, ok := allowedColumns[row]; !ok {
        return errors.New("invalid column name")
    }

    query := fmt.Sprintf("UPDATE pro_users SET %s=$1 WHERE id=$2", row)

    _, err := r.DB.Exec(query, value, id)

    if err != nil {
        log.Printf("Could not update Pro User row name: %v", row)
        return err
    }

    return nil
}

func (r *pGProUserRepository) FindByID(ctx context.Context, id int) (*model.ProUser, error) {
    proUser := &model.ProUser{}

    query := "SELECT * FROM pro_users WHERE id=$1"

    if err := r.DB.GetContext(ctx, proUser, query, id); err != nil {
        return nil, err
    }

    return proUser, nil
}

func (r *pGProUserRepository) FindByEmail(ctx context.Context, email string) (*model.ProUser, error) {
    proUser := &model.ProUser{}

    query := "SELECT * FROM pro_users WHERE email=$1"

    if err := r.DB.GetContext(ctx, proUser, query, email); err != nil {
        log.Printf("Unable to get Pro User with email: %v. Err: %v\n", email, err)
        return nil, err
    }

    return proUser, nil
}

func (r *pGProUserRepository) UpdateInitials(ctx context.Context, p *model.ProUser) error {
    query := `UPDATE pro_users SET name=$1, surname=$2, patronym=$3 WHERE id=$4`

    _, err := r.DB.ExecContext(ctx, query, p.Name, p.Surname, p.Patronym, p.Id)

    if err != nil {
        return err
    }

    return nil
}

func (r *pGProUserRepository) AddShop(ctx context.Context, shops []model.Shop, id int) error {
    for _, shop := range shops {
        addressJson, err := json.Marshal(shop.Addresses)

        if err != nil {
            log.Printf("Unable to marshal %v", err)
            return err
        }

        workScheduleJson, err := json.Marshal(shop.WorkSchedule)

        if err != nil {
            log.Printf("Error marshaling work schedule: %v", err)
            return err
        }

        query := `INSERT INTO shop (name, description, addresses, shop_logo, work_schedule, pro_user_id)
                    VALUES ($1, $2, $3, $4, $5, $6)`

        _, err = r.DB.ExecContext(ctx, query, shop.Name, shop.Description, addressJson, shop.ShopLogo, workScheduleJson, id)

        if err != nil {
            log.Printf("Unable to write shop %v", err)
            return err
        }
    }

    return nil
}

func (r *pGProUserRepository) GetOnboardingPath(ctx context.Context, id int) (string, error) {
    var proUser model.ProUser

    query := `SELECT * from pro_users WHERE id=$1`

    err := r.DB.GetContext(ctx, &proUser, query, id)

    if err != nil {
        return "", err
    }

    if proUser.Role == "individual" {
        if proUser.City == nil {
            return "city", nil
        } else if proUser.Name == nil {
            return "initials", nil
        } else {
            return "done", nil
        }
    } else if proUser.Role == "shop" {
        var byteAddresses []byte
        var addresses *[]model.Address
        query := `SELECT addresses FROM shop WHERE pro_user_id=$1`

        err := r.DB.GetContext(ctx, &byteAddresses, query, id)

        if err != nil {
            if err == sql.ErrNoRows {
                return "address", nil
            }
            log.Printf("Unable to get shop %v", err)
            return "", err
        }

        err = json.Unmarshal(byteAddresses, &addresses)

        if err != nil {
            log.Printf("Unable to unmarshal addresses %v", err)
            return "", err
        }

        if addresses == nil {
            return "address", nil
        } else if len(*addresses) == 0 {
            return "address", nil
        } else {
            return "done", nil
        }
    } else {
        return "role", err
    }

}

func (r *pGProUserRepository) CreateOrder(ctx context.Context, orders []model.Order, id int) error {
    for _, order := range orders {

        query := `INSERT INTO orders (pro_user_id, plantation_id, delivery_id, delivery_address, quantity, price_for_one, tenge_price_for_one, total_price, total_tenge_price)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

        walletQuery := `SELECT wallet FROM pro_users WHERE id=$1`
        var wallet float64

        if err := r.DB.GetContext(ctx, &wallet, walletQuery, id); err != nil {
            log.Printf("Error getting wallet: %v", err)
            return err
        }

        if wallet < order.TotalPrice {
            err := errors.New("not enough money on wallet")
            log.Printf("Not enough money on wallet, money: %v", wallet)
            return err
        }

        _, err := r.DB.ExecContext(ctx, query, id, order.PlantationId, order.DeliveryId, order.DeliveryAddress, order.Quantity, order.PriceForOne, order.TengePriceForOne, order.TotalPrice, order.TotalTengePrice)

        withdrawWalletQuery := `UPDATE pro_users SET wallet= wallet - $1 WHERE id=$2`

        _, err = r.DB.ExecContext(ctx, withdrawWalletQuery, order.TotalTengePrice, id)

        if err != nil {
            log.Printf("Error withdrawing wallet: %v", err)
            return err
        }

        if err != nil {
            log.Printf("Unable to create order %v", err)
            return err
        }
    }

    return nil

}

func (r *pGProUserRepository) GetCardInfo(ctx context.Context, deliveryId int) (*model.Delivery, *int, *[]model.Address, error) {
    var delivery model.Delivery
    var residue int
    var addresses []model.Address

    deliveryQuery := `SELECT * FROM delivery WHERE id=$1`

    if err := r.DB.GetContext(ctx, &delivery, deliveryQuery, deliveryId); err != nil {
        log.Printf("Error fetching delivery %v", err)
        return nil, nil, nil, err
    }

    residueQuery := `SELECT COUNT(*) FROM delivery WHERE length=$1 AND color=$2 AND species=$3 AND box_size=$4 AND product=$5`

    if err := r.DB.GetContext(ctx, &residue, residueQuery, delivery.Length, delivery.Color, delivery.Species, delivery.BoxSize, delivery.Product); err != nil {
        log.Printf("Error fetching residue %v", err)
        return nil, nil, nil, err
    }

    addressesQuery := `SELECT city, street, house FROM storehouses`

    if err := r.DB.SelectContext(ctx, &addresses, addressesQuery); err != nil {
        log.Printf("Error fetching addresses from storehouse %v", err)
        return nil, nil, nil, err
    }

    return &delivery, &residue, &addresses, nil

}

func (r *pGProUserRepository) GetOrders(ctx context.Context, id int, status string) (*[]model.ProUserOrder, error) {
    var orders []model.ProUserOrder
    query := `SELECT d.species, d.box_size, o.id, p.name, o.plantation_id, o.price_for_one, o.tenge_price_for_one, o.total_tenge_price, o.delivery_address, o.status FROM orders o
        JOIN delivery d on d.id = o.delivery_id JOIN plantations p on p.id = d.plantation_id WHERE o.pro_user_id=$1 AND o.status=$2`

    if err := r.DB.SelectContext(ctx, &orders, query, id, status); err != nil {
        log.Printf("Couldn't fetch pro user orders %v", err)
        return nil, err
    }

    return &orders, nil
}

func (r *pGProUserRepository) GetFinancialCharges(ctx context.Context) (*model.FinancialCharges, error) {
    var charges model.FinancialCharges

    chargesQuery := `SELECT * FROM financial_charges`

    if err := r.DB.GetContext(ctx, &charges, chargesQuery); err != nil {
        log.Printf("Error fetching financial charges %v", err)
        return nil, err
    }

    return &charges, nil
}

func (r *pGProUserRepository) GetWallet(ctx context.Context, id int) (float64, error) {
    var wallet float64
    query := `SELECT wallet FROM pro_users WHERE id=$1`

    if err := r.DB.GetContext(ctx, &wallet, query, id); err != nil {
        log.Printf("Error getting wallet: %v", err)
        return 0, err
    }

    return wallet, nil
}

func (r *pGProUserRepository) GetProUserInfo(ctx context.Context, id int) (*model.PrivateOffice, error) {
    var proUserInfo model.ProUserInfo
    var shops []model.Shop

    proUserQuery := `SELECT name, surname, patronym, city, role FROM pro_users WHERE id=$1`

    if err := r.DB.GetContext(ctx, &proUserInfo, proUserQuery, id); err != nil {
        log.Printf("Error getting pro user info: %v", err)
        return nil, err
    }

    shopsQuery := `SELECT * FROM shop WHERE pro_user_id=$1`

    err := r.DB.SelectContext(ctx, &shops, shopsQuery, id)

    if err != nil {
        log.Printf("Error getting shops: %v", err)
        return nil, err
    }

    info := &model.PrivateOffice{
        ProUserInfo: proUserInfo,
        Shops:       shops,
    }

    return info, nil
}

func (r *pGProUserRepository) UpdateShopsInfo(ctx context.Context, shops []model.Shop, id int) error {
    for _, shop := range shops {
        if shop.Id != 0 {
            if err := updateShop(ctx, r.DB, &shop, id); err != nil {
                log.Printf("Error updating shop: %v", err)
                return err
            }
        } else {
            if err := addShop(ctx, r.DB, shop, id); err != nil {
                log.Printf("Error adding shop: %v", err)
                return err
            }
        }
    }

    return nil
}

func (r *pGProUserRepository) AddFavoriteProduct(ctx context.Context, proUserId, deliveryId int) error {
    query := `INSERT INTO favorite_delivery (pro_user_id, delivery_id) VALUES ($1, $2)`

    _, err := r.DB.ExecContext(ctx, query, proUserId, deliveryId)

    if err != nil {
        log.Printf("Error inserting into favorite delivery table: %v", err)
        return err
    }

    return nil
}

func (r *pGProUserRepository) RemoveFavoriteProduct(ctx context.Context, proUserId, deliveryId int) error {
    log.Printf("tytytyty %v,%v", proUserId, deliveryId)
    query := `DELETE FROM favorite_delivery WHERE pro_user_id=$1 AND delivery_id=$2`

    _, err := r.DB.ExecContext(ctx, query, proUserId, deliveryId)

    if err != nil {
        log.Printf("Error deleting from favorite_delivery table: %v", err)
        return err
    }

    return nil
}

func (r *pGProUserRepository) GetFavoriteProducts(ctx context.Context, id int) ([]model.Deliveries, error) {
    var deliveries []model.Deliveries
    query := `SELECT d.*, p.name AS plantation_name FROM delivery d JOIN favorite_delivery fd on d.id = fd.delivery_id JOIN plantations p on p.id = d.plantation_id WHERE fd.pro_user_id=$1`

    rows, err := r.DB.QueryContext(ctx, query, id)

    if err != nil {
        log.Printf("Error getting products: %v", err)
        return nil, err
    }

    for rows.Next() {
        var delivery model.Deliveries
        if err := rows.Scan(
            &delivery.Id, &delivery.FarmBox, &delivery.BoxSize, &delivery.Mixed,
            &delivery.Species, &delivery.Product, &delivery.Color, &delivery.Length,
            &delivery.Price, &delivery.Boxes, &delivery.Packing, &delivery.PlantationId,
            &delivery.PlantationName); err != nil {
            log.Printf("Error scanning row: %v", err)
            return nil, err
        }
        deliveries = append(deliveries, delivery)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows: %v", err)
        return nil, err
    }

    return deliveries, nil

}

func (r *pGProUserRepository) UpdateProUserRole(ctx context.Context, id int, role string) error {
    query := `UPDATE pro_users SET role=$1 WHERE id=$2`

    _, err := r.DB.ExecContext(ctx, query, role, id)

    if err != nil {
        log.Printf("Error Updating pro user role: %v", err)
        return err
    }

    return nil
}

func (r *pGProUserRepository) GetProUserRole(ctx context.Context, id int) (*string, error) {
    var role *string

    query := `SELECT role FROM pro_users WHERE id=$1`

    if err := r.DB.GetContext(ctx, &role, query, id); err != nil {
        log.Printf("Error getting pro user role %v", err)
        return nil, err
    }

    return role, nil
}

func (r *pGProUserRepository) AddToCart(ctx context.Context, cart model.Cart) error {
    query := `INSERT INTO cart (pro_user_id, delivery_id, delivery_address, quantity, total_price) VALUES ($1, $2, $3, $4, $5)`

    _, err := r.DB.ExecContext(ctx, query, cart.ProUserId, cart.DeliveryId, cart.DeliveryAddress, cart.Quantity, cart.TotalPrice)

    if err != nil {
        log.Printf("Error adding item to cart: %v", err)
        return err
    }

    return nil
}

func (r *pGProUserRepository) GetCart(ctx context.Context, id int) ([]model.CartItem, error) {
    var cartItems []model.CartItem

    query := `SELECT d.*, c.delivery_address, c.quantity, c.status, c.created_at, c.time_to_overdue FROM cart c JOIN delivery d on d.id = c.delivery_id WHERE c.pro_user_id=$1`

    err := r.DB.SelectContext(ctx, &cartItems, query, id)

    if err != nil {
        log.Printf("Error getting cart items: %v", err)
    }

    for i := range cartItems {
        cartItems[i].TotalPrice = cartItems[i].Price * float64(cartItems[i].Quantity)
    }

    return cartItems, nil
}

func (r *pGProUserRepository) RemoveFromCart(ctx context.Context, proUserId, deliveryId int) error {
    query := `DELETE FROM favorite_delivery WHERE pro_user_id=$1 AND delivery_id=$2`

    _, err := r.DB.ExecContext(ctx, query, proUserId, deliveryId)

    if err != nil {
        log.Printf("Error deleting from cart: %v", err)
        return err
    }

    return nil
}

func updateShop(ctx context.Context, d *sqlx.DB, shopInfo *model.Shop, id int) error {
    if len(shopInfo.Name) != 0 {
        _, err := d.ExecContext(ctx, "UPDATE shop SET name = $1 WHERE id = $2", shopInfo.Name, id)
        if err != nil {
            return err
        }
    }
    if len(shopInfo.Description) != 0 {
        _, err := d.ExecContext(ctx, "UPDATE shop SET description = $1 WHERE id = $2", shopInfo.Description, id)
        if err != nil {
            return err
        }
    }
    if shopInfo.Addresses != nil {
        _, err := d.ExecContext(ctx, "UPDATE shop SET addresses = $1 WHERE id = $2", shopInfo.Addresses, id)
        if err != nil {
            return err
        }
    }

    if shopInfo.ShopLogo != nil {
        _, err := d.ExecContext(ctx, "UPDATE shop SET shop_logo = $1 WHERE id = $2", *shopInfo.ShopLogo, id)
        if err != nil {
            return err
        }
    }
    if shopInfo.WorkSchedule.Days != nil {
        workScheduleJson, err := json.Marshal(shopInfo.WorkSchedule)

        if err != nil {
            log.Printf("Error marshaling work schedule: %v", err)
            return err
        }
        _, err = d.ExecContext(ctx, "UPDATE shop SET work_schedule = $1 WHERE id = $2", workScheduleJson, id)
        if err != nil {
            return err
        }
    }
    return nil
}

func addShop(ctx context.Context, d *sqlx.DB, shop model.Shop, id int) error {

    workScheduleJson, err := json.Marshal(shop.WorkSchedule)

    if err != nil {
        log.Printf("Error marshaling work schedule: %v", err)
        return err
    }
    _, err = d.ExecContext(ctx, `INSERT INTO shop (name, description, addresses, shop_logo, work_schedule, pro_user_id) VALUES ($1, $2, $3, $4, $5, $6)`, shop.Name, shop.Description, shop.Addresses, shop.ShopLogo, workScheduleJson, id)
    if err != nil {
        log.Printf("Error adding shop %v", err)
        return err
    }

    return nil
}
