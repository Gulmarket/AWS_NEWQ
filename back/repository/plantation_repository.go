package repository

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gulmarket/model"
    "log"

    "github.com/jmoiron/sqlx"
    "github.com/lib/pq"
)

type pGPlantationRepository struct {
    DB *sqlx.DB
}

func NewPlantationRepository(db *sqlx.DB) model.PlantationRepository {
    return &pGPlantationRepository{
        DB: db,
    }
}

func (r *pGPlantationRepository) Create(ctx context.Context, p *model.Plantation) error {
    query := `INSERT INTO plantations (email, password, phone)
              VALUES ($1, $2, $3) RETURNING *`

    if err := r.DB.GetContext(ctx, p, query, p.Email, p.Password, p.Phone); err != nil {
        if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
            log.Printf("Could not create a plantation with email: %v. Reason: %v\n", p.Email, err.Code.Name())
            return err
        }

        log.Printf("Could not create a plantation with email: %v. Reason: %v\n", p.Email, err)
        return err
    }
    return nil
}

func (r *pGPlantationRepository) UpdatePlantationRow(ctx context.Context, p *model.Plantation, row string) error {
    allowedColumns := map[string]bool{
        "email":    true,
        "password": true,
        "name":     true,
        "phone":    true,
    }

    if _, ok := allowedColumns[row]; !ok {
        return errors.New("invalid column name")
    }

    query := fmt.Sprintf("UPDATE plantations SET %s=$1 WHERE id=$2", row)

    _, err := r.DB.Exec(query, p.Password, p.Id)

    if err != nil {
        log.Printf("Could not update plantation password")
        return err
    }

    return nil
}

func (r *pGPlantationRepository) UpdateName(ctx context.Context, p *model.Plantation) error {
    query := `UPDATE plantations SET name = $2 WHERE id = $1`

    _, err := r.DB.ExecContext(ctx, query, p.Id, p.Name)

    if err != nil {
        return err
    }

    return nil
}

func (r *pGPlantationRepository) FindByID(ctx context.Context, id int) (*model.Plantation, error) {
    plantation := &model.Plantation{}

    query := "SELECT * FROM plantations WHERE id=$1"

    if err := r.DB.GetContext(ctx, plantation, query, id); err != nil {
        return nil, err
    }

    return plantation, nil
}

func (r *pGPlantationRepository) FindByEmail(ctx context.Context, email string) (*model.Plantation, error) {
    plantation := &model.Plantation{}

    query := "SELECT * FROM plantations WHERE email=$1"

    if err := r.DB.GetContext(ctx, plantation, query, email); err != nil {
        log.Printf("Unable to get plantation with email: %v. Err: %v\n", email, err)
        return nil, err
    }

    return plantation, nil
}

func (r *pGPlantationRepository) UpdateLocation(ctx context.Context, p *model.Plantation, id int) error {
    query := `UPDATE plantations SET country=$1, city=$2 WHERE id=$3`

    _, err := r.DB.ExecContext(ctx, query, p.Country, p.City, id)

    if err != nil {
        log.Printf("Failed to update plantation location: %v", err)
        return err
    }

    return nil
}

func (r *pGPlantationRepository) UpdateInfo(ctx context.Context, p *model.Plantation, id int) error {
    query := `UPDATE plantations SET country=$1, city=$2, name=$3, description=$4, logo_url=$5, work_schedule=$6 WHERE id=$7`

    workScheduleJson, err := json.Marshal(p.WorkSchedule)

    if err != nil {
        log.Printf("Error marshaling work schedule: %v", err)
        return err
    }

    _, err = r.DB.ExecContext(ctx, query, p.Country, p.City, p.Name, p.Description, p.LogoUrl, workScheduleJson, id)

    if err != nil {
        log.Printf("Failed to update plantation info: %v", err)
        return err
    }

    return nil
}

func (r *pGPlantationRepository) GetAllDelivery(ctx context.Context) (*[]model.Delivery, error) {
    delivery := &[]model.Delivery{}

    query := "SELECT * FROM delivery"

    if err := r.DB.GetContext(ctx, delivery, query); err != nil {
        log.Printf("Unable to get delivery: %v. Err: %v\n", err)
        return nil, err
    }

    return delivery, nil
}

func (r *pGPlantationRepository) GetFilteredProducts(ctx context.Context, params model.Filter, proUserId int) (*model.FlowersFilter, error) {
    var flowerTypes []string
    var deliveries []model.Deliveries
    var args []interface{}
    query := `
        SELECT 
            d.id, d.farm_box, d.box_size, d.mixed, d.species, d.product, d.color,
            d.length, d.price, d.boxes, d.packing, d.plantation_id, p.name AS plantation_name,
            CASE WHEN fd.delivery_id IS NOT NULL THEN true ELSE false END AS is_favorite
        FROM 
            delivery d
        JOIN 
            plantations p ON d.plantation_id = p.id
        LEFT JOIN 
            favorite_delivery fd ON d.id = fd.delivery_id AND fd.pro_user_id = $1
        WHERE 
            CAST(regexp_replace(d.length, '[^0-9]', '', 'g') AS INTEGER) BETWEEN $2 AND $3
    `

    args = append(args, proUserId, params.SizeStart, params.SizeEnd)

    argIndex := 4

    if len(params.FlowerType) != 0 {
        query += fmt.Sprintf(" AND d.species = ANY($%d)", argIndex)
        args = append(args, pq.Array(params.FlowerType))
        argIndex++
    }

    if len(params.FlowerSpecies) != 0 {
        query += fmt.Sprintf(" AND d.product = ANY($%d)", argIndex)
        args = append(args, pq.Array(params.FlowerSpecies))
        argIndex++
    }

    if len(params.BoxType) != 0 {
        query += fmt.Sprintf(" AND d.box_size = ANY($%d)", argIndex)
        args = append(args, pq.Array(params.BoxType))
        argIndex++
    }

    if len(params.Color) != 0 {
        query += fmt.Sprintf(" AND d.color = ANY($%d)", argIndex)
        args = append(args, pq.Array(params.Color))
    }

    rows, err := r.DB.QueryContext(ctx, query, args...)
    if err != nil {
        log.Printf("Error getting products: %v", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var delivery model.Deliveries
        if err := rows.Scan(
            &delivery.Id, &delivery.FarmBox, &delivery.BoxSize, &delivery.Mixed,
            &delivery.Species, &delivery.Product, &delivery.Color, &delivery.Length,
            &delivery.Price, &delivery.Boxes, &delivery.Packing, &delivery.PlantationId,
            &delivery.PlantationName, &delivery.IsFavorite); err != nil {
            log.Printf("Error scanning row: %v", err)
            return nil, err
        }
        deliveries = append(deliveries, delivery)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows: %v", err)
        return nil, err
    }

    flowerTypesQuery := `SELECT DISTINCT species FROM delivery`

    err = r.DB.SelectContext(ctx, &flowerTypes, flowerTypesQuery)
    if err != nil {
        log.Printf("Failed to get flower types: %v", err)
        return nil, err
    }

    flowerSpeciesMap := make(map[string][]string)
    if len(params.FlowerType) != 0 {
        flowerSpeciesQuery := `SELECT species, product FROM delivery WHERE species = ANY($1) GROUP BY species, product`
        rows, err := r.DB.QueryContext(ctx, flowerSpeciesQuery, pq.Array(params.FlowerType))
        if err != nil {
            log.Printf("Failed to get flower species: %v", err)
            return nil, err
        }
        defer rows.Close()

        for rows.Next() {
            var species, product string
            if err := rows.Scan(&species, &product); err != nil {
                log.Printf("Error scanning row: %v", err)
                return nil, err
            }
            flowerSpeciesMap[species] = append(flowerSpeciesMap[species], product)
        }
        if err := rows.Err(); err != nil {
            log.Printf("Error iterating rows: %v", err)
            return nil, err
        }
    }

    result := &model.FlowersFilter{
        FlowerTypes:   flowerTypes,
        FlowerSpecies: flowerSpeciesMap,
        Deliveries:    deliveries,
    }

    return result, nil
}

func (r *pGPlantationRepository) GetOrders(ctx context.Context, id int, status string) (*[]model.PlantationOrder, error) {
    var orders []model.PlantationOrder
    query := `SELECT d.species, d.box_size, o.id, o.total_price, o.delivery_address, o.status, o.quantity FROM orders o
        JOIN delivery d on d.id = o.delivery_id JOIN pro_users p on p.id = d.plantation_id WHERE o.plantation_id=$1 AND o.status=$2`

    if err := r.DB.SelectContext(ctx, &orders, query, id, status); err != nil {
        log.Printf("Couldn't fetch plantations orders %v", err)
        return nil, err
    }

    return &orders, nil
}

func (r *pGPlantationRepository) DeclineOrder(ctx context.Context, orderId int) error {
    query := `UPDATE orders SET status='plantation_declined' WHERE id=$1`

    _, err := r.DB.ExecContext(ctx, query, orderId)

    if err != nil {
        log.Printf("Unable to decline plantation order: %v", err)
        return err
    }

    return nil
}

func (r *pGPlantationRepository) ApproveOrder(ctx context.Context, orderId int) error {
    query := `UPDATE orders SET status='plantation_approved' WHERE id=$1`

    _, err := r.DB.ExecContext(ctx, query, orderId)

    if err != nil {
        log.Printf("Unable to approve plantation order: %v", err)
        return err
    }

    return nil
}

func (r *pGPlantationRepository) GetPlantationFilteredProducts(ctx context.Context, params model.Filter, id int) (*model.PlantationFlowersFilter, error) {
    var flowerTypes []string
    var deliveries []model.Delivery
    var args []interface{}
    query := `
        SELECT 
            d.id, d.farm_box, d.box_size, d.mixed, d.species, d.product, d.color,
            d.length, d.price, d.boxes, d.packing, d.plantation_id
        FROM 
            delivery d
        WHERE 
            CAST(regexp_replace(d.length, '[^0-9]', '', 'g') AS INTEGER) BETWEEN $1 AND $2 AND plantation_id=$3
    `

    args = append(args, params.SizeStart, params.SizeEnd, id)

    argIndex := 4

    if len(params.FlowerType) != 0 {
        query += fmt.Sprintf(" AND d.species = ANY($%d)", argIndex)
        args = append(args, pq.Array(params.FlowerType))
        argIndex++
    }

    if len(params.FlowerSpecies) != 0 {
        query += fmt.Sprintf(" AND d.product = ANY($%d)", argIndex)
        args = append(args, pq.Array(params.FlowerSpecies))
        argIndex++
    }

    if len(params.BoxType) != 0 {
        query += fmt.Sprintf(" AND d.box_size = ANY($%d)", argIndex)
        args = append(args, pq.Array(params.BoxType))
        argIndex++
    }

    if len(params.Color) != 0 {
        query += fmt.Sprintf(" AND d.color = ANY($%d)", argIndex)
        args = append(args, pq.Array(params.Color))
    }

    rows, err := r.DB.QueryContext(ctx, query, args...)
    if err != nil {
        log.Printf("Error getting products: %v", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var delivery model.Delivery
        if err := rows.Scan(
            &delivery.Id, &delivery.FarmBox, &delivery.BoxSize, &delivery.Mixed,
            &delivery.Species, &delivery.Product, &delivery.Color, &delivery.Length,
            &delivery.Price, &delivery.Boxes, &delivery.Packing, &delivery.PlantationId); err != nil {
            log.Printf("Error scanning row: %v", err)
            return nil, err
        }
        deliveries = append(deliveries, delivery)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows: %v", err)
        return nil, err
    }

    flowerTypesQuery := `SELECT DISTINCT species FROM delivery`

    err = r.DB.SelectContext(ctx, &flowerTypes, flowerTypesQuery)
    if err != nil {
        log.Printf("Failed to get flower types: %v", err)
        return nil, err
    }

    flowerSpeciesMap := make(map[string][]string)
    if len(params.FlowerType) != 0 {
        flowerSpeciesQuery := `SELECT species, product FROM delivery WHERE species = ANY($1) GROUP BY species, product`
        rows, err := r.DB.QueryContext(ctx, flowerSpeciesQuery, pq.Array(params.FlowerType))
        if err != nil {
            log.Printf("Failed to get flower species: %v", err)
            return nil, err
        }
        defer rows.Close()

        for rows.Next() {
            var species, product string
            if err := rows.Scan(&species, &product); err != nil {
                log.Printf("Error scanning row: %v", err)
                return nil, err
            }
            flowerSpeciesMap[species] = append(flowerSpeciesMap[species], product)
        }
        if err := rows.Err(); err != nil {
            log.Printf("Error iterating rows: %v", err)
            return nil, err
        }
    }

    result := &model.PlantationFlowersFilter{
        FlowerTypes:   flowerTypes,
        FlowerSpecies: flowerSpeciesMap,
        Deliveries:    deliveries,
    }

    return result, nil
}
