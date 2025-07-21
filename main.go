package main

import (
	"fmt"
	"log"
)

func main() {
	// Ejemplos de JSON para probar
	examples := []string{
		`{"name": "Juan", "age": 30, "active": true}`,
		`[1, 2, 3, "hello", true, null]`,
		`{
			"usuario": {
				"nombre": "María",
				"datos": {
					"edad": 25,
					"emails": ["maria@example.com", "m.garcia@work.com"]
				}
			},
			"configuracion": {
				"tema": "dark",
				"notificaciones": false
			}
		}`,
		`[]`,
		`{}`,
		`null`,
		`"simple string"`,
		`42.5`,
		`true`,
	}

	for i, jsonStr := range examples {
		fmt.Printf("=== Ejemplo %d ===\n", i+1)
		fmt.Printf("JSON de entrada: %s\n", jsonStr)

		result, err := ParseJSON(jsonStr)
		if err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Resultado parseado: %+v\n", result)
			fmt.Printf("Tipo: %T\n", result)
		}
		fmt.Println()
	}

	// Ejemplos de JSON inválido para probar el manejo de errores
	invalidExamples := []string{
		`{"name": }`,            // valor faltante
		`{"name" "value"}`,      // ':' faltante
		`[1, 2, 3,]`,            // coma final
		`{"unclosed": "string}`, // string sin cerrar
		`{name: "value"}`,       // clave sin comillas
		`[1, 2 3]`,              // coma faltante
		`truee`,                 // booleano inválido
		`{"a": 1 "b": 2}`,       // coma faltante entre pares
	}

	fmt.Println("=== Pruebas de JSON inválido ===")
	for i, jsonStr := range invalidExamples {
		fmt.Printf("Prueba %d: %s\n", i+1, jsonStr)
		_, err := ParseJSON(jsonStr)
		if err != nil {
			fmt.Printf("Error (esperado): %v\n", err)
		} else {
			fmt.Printf("¡Inesperado! No se detectó error\n")
		}
		fmt.Println()
	}
}
