package pagos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Cliente estructura simple que coincide con la tabla
type Cliente struct {
	ID                  uint       `json:"id" gorm:"primaryKey;column:ID"`
	Nombre              string     `json:"nombre" gorm:"column:NOMBRE"`
	Apellido            string     `json:"apellido" gorm:"column:APELLIDO"`
	DUI                 string     `json:"dui" gorm:"column:DUI"`
	Direccion           string     `json:"direccion" gorm:"column:DIRECCION"`
	Telefono            string     `json:"telefono" gorm:"column:TELEFONO"`
	Celular             string     `json:"celular" gorm:"column:CELULAR"`
	FechaIngreso        *time.Time `json:"fecha_ingreso" gorm:"column:FECHA_INGRESO"`
	Departamento        string     `json:"departamento" gorm:"column:DEPARTAMENTO"`
	Activo              *int8      `json:"activo" gorm:"column:ACTIVO"`
	Giro                string     `json:"giro" gorm:"column:GIRO"`
	Referencia1         string     `json:"referencia1" gorm:"column:REFERENCIA1"`
	TelRef1             string     `json:"tel_ref1" gorm:"column:TELREF1"`
	Referencia2         string     `json:"referencia2" gorm:"column:REFERENCIA2"`
	TelRef2             string     `json:"tel_ref2" gorm:"column:TELREF2"`
	IDGestor            *int       `json:"id_gestor" gorm:"column:IDGESTOR"`
	TipoPer             string     `json:"tipo_per" gorm:"column:TIPO_PER"`
	FechaNacimiento     *time.Time `json:"fecha_nacimiento" gorm:"column:FECHA_NACIMIENTO"`
	NIT                 string     `json:"nit" gorm:"column:NIT"`
	Sexo                string     `json:"sexo" gorm:"column:SEXO"`
	DUIFrente           string     `json:"dui_frente" gorm:"column:DUI_FRENTE"`
	DUIDetras           string     `json:"dui_detras" gorm:"column:DUI_DETRAS"`
	FotoNegocio1        string     `json:"foto_negocio1" gorm:"column:FOTONEGOCIO1"`
	FotoNegocio2        string     `json:"foto_negocio2" gorm:"column:FOTONEGOCIO2"`
	FotoNegocio3        string     `json:"foto_negocio3" gorm:"column:FOTONEGOCIO3"`
	FotoNegocio4        string     `json:"foto_negocio4" gorm:"column:FOTONEGOCIO4"`
	Longitud            *float64   `json:"longitud" gorm:"column:LONGITUD"`
	Latitud             *float64   `json:"latitud" gorm:"column:LATITUD"`
	Profesion           string     `json:"profesion" gorm:"column:PROFESION"`
	Email               string     `json:"email" gorm:"column:EMAIL"`
}

func (Cliente) TableName() string {
	return "cliente"
}

// ClienteController maneja las operaciones de clientes
type ClienteController struct {
	db *gorm.DB
}

// getEnvOrDefault obtiene una variable de entorno o devuelve un valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewClienteController crea una nueva instancia del controlador
func NewClienteController() *ClienteController {
	// Obtener configuración desde variables de entorno (SIN valores por defecto sensibles)
	host := os.Getenv("PAGOS_DB_HOST")
	port := os.Getenv("PAGOS_DB_PORT")
	username := os.Getenv("PAGOS_DB_USERNAME")
	password := os.Getenv("PAGOS_DB_PASSWORD")
	database := os.Getenv("PAGOS_DB_DATABASE")
	charset := getEnvOrDefault("PAGOS_DB_CHARSET", "utf8mb4")
	
	// Validar que las variables requeridas existan
	if host == "" || port == "" || username == "" || password == "" || database == "" {
		panic("Missing required database environment variables: PAGOS_DB_HOST, PAGOS_DB_PORT, PAGOS_DB_USERNAME, PAGOS_DB_PASSWORD, PAGOS_DB_DATABASE")
	}
	
	// Construir DSN desde variables de entorno
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&tls=required",
		username, password, host, port, database, charset)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// Log del error sin exponer credenciales
		panic(fmt.Sprintf("Failed to connect to pagos database at %s:%s - %v", host, port, err))
	}

	return &ClienteController{db: db}
}

// GetAllClientes obtiene todos los clientes con paginación
// @Summary Obtener todos los clientes
// @Description Obtiene una lista paginada de todos los clientes
// @Tags Clientes
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} map[string]interface{} "Lista de clientes"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /api/pagos/clientes [get]
func (c *ClienteController) GetAllClientes(w http.ResponseWriter, r *http.Request) {
	// Parámetros de paginación
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var clientes []Cliente
	var total int64

	// Contar total
	c.db.Model(&Cliente{}).Count(&total)

	// Obtener clientes con paginación
	offset := (page - 1) * limit
	result := c.db.Offset(offset).Limit(limit).Order("ID DESC").Find(&clientes)

	if result.Error != nil {
		http.Error(w, "Error al obtener clientes", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":        clientes,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetClienteByID obtiene un cliente por ID
// @Summary Obtener cliente por ID
// @Description Obtiene un cliente específico por su ID
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path int true "ID del cliente"
// @Success 200 {object} Cliente "Cliente encontrado"
// @Failure 404 {object} map[string]interface{} "Cliente no encontrado"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /api/pagos/clientes/{id} [get]
func (c *ClienteController) GetClienteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var cliente Cliente
	result := c.db.Where("ID = ?", id).First(&cliente)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Cliente no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener cliente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

// GetClienteByDUI obtiene un cliente por DUI
// @Summary Obtener cliente por DUI
// @Description Obtiene un cliente específico por su DUI
// @Tags Clientes
// @Accept json
// @Produce json
// @Param dui path string true "DUI del cliente"
// @Success 200 {object} Cliente "Cliente encontrado"
// @Failure 404 {object} map[string]interface{} "Cliente no encontrado"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /api/pagos/clientes/dui/{dui} [get]
func (c *ClienteController) GetClienteByDUI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dui := vars["dui"]

	var cliente Cliente
	result := c.db.Where("DUI = ?", dui).First(&cliente)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Cliente no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener cliente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

// SearchClientes busca clientes por nombre, apellido o DUI
// @Summary Buscar clientes
// @Description Busca clientes por nombre, apellido o DUI
// @Tags Clientes
// @Accept json
// @Produce json
// @Param q query string true "Término de búsqueda"
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} map[string]interface{} "Resultados de búsqueda"
// @Failure 400 {object} map[string]interface{} "Parámetros inválidos"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /api/pagos/clientes/search [get]
func (c *ClienteController) SearchClientes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Parámetro de búsqueda 'q' es requerido", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var clientes []Cliente
	var total int64

	searchTerm := fmt.Sprintf("%%%s%%", query)
	whereClause := c.db.Where("NOMBRE LIKE ? OR APELLIDO LIKE ? OR DUI LIKE ? OR CONCAT(NOMBRE, ' ', APELLIDO) LIKE ?",
		searchTerm, searchTerm, searchTerm, searchTerm)
	
	// Contar total de resultados
	whereClause.Model(&Cliente{}).Count(&total)
	
	// Obtener clientes con paginación
	offset := (page - 1) * limit
	result := whereClause.Offset(offset).Limit(limit).Order("ID DESC").Find(&clientes)

	if result.Error != nil {
		http.Error(w, "Error al buscar clientes", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":        clientes,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
		"search_term": query,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetActiveClientes obtiene solo los clientes activos
// @Summary Obtener clientes activos
// @Description Obtiene una lista paginada de clientes con estado activo
// @Tags Clientes
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Success 200 {object} map[string]interface{} "Lista de clientes activos"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /api/pagos/clientes/activos [get]
func (c *ClienteController) GetActiveClientes(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var clientes []Cliente
	var total int64

	// Contar total de clientes activos
	c.db.Model(&Cliente{}).Where("ACTIVO = ?", 1).Count(&total)

	// Obtener clientes activos con paginación
	offset := (page - 1) * limit
	result := c.db.Where("ACTIVO = ?", 1).Offset(offset).Limit(limit).Order("ID DESC").Find(&clientes)

	if result.Error != nil {
		http.Error(w, "Error al obtener clientes activos", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":        clientes,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
		"filter":      "active_only",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}