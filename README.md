# Reto 02: Parser JSON en Go

Este proyecto consiste en la implementación de un **parser JSON de bajo nivel** utilizando el lenguaje de programación Go. El parser ha sido desarrollado **sin utilizar el paquete `encoding/json`**, cumpliendo así con las restricciones del reto.

## 📁 Estructura del Proyecto
📁 RETO02_PARSER_JSON/
├── 📄 go.mod
├── 📄 main.go // Código principal con ejemplos de uso del parser
├── 📄 parser.go // Implementación completa del parser JSON
├── 📄 parser_test.go // Casos de prueba unitarios
└── 📄 README.md

## 🧠 ¿Qué hace este parser?

Este parser convierte cadenas JSON válidas en estructuras nativas de Go:

- Objetos JSON → `map[string]interface{}`
- Arreglos JSON → `[]interface{}`
- Strings → `string`
- Números → `float64`
- Booleanos → `bool`
- Null → `nil`

Además, **detecta errores** comunes en JSON mal formados, indicando la línea y columna del error.

## ⚙️ Compilación y Ejecución

Desde la carpeta raíz del proyecto:

```bash
go run .

Para ejecutar las pruebas unitarias:
go test

✅ Características implementadas
✔️ Soporte para objetos JSON anidados

✔️ Soporte para arreglos anidados

✔️ Manejo de strings con comillas escapadas

✔️ Reconocimiento de números decimales y enteros

✔️ Soporte completo para true, false y null

✔️ Reporte detallado de errores con línea y columna

✔️ Pruebas para casos válidos e inválidos

📋 Casos de prueba cubiertos
JSON válidos:

Objetos simples y anidados

Arreglos con diferentes tipos

Strings simples y con escapado

Booleanos y null

Valores primitivos como "texto", 42, true

JSON inválidos:

Claves sin comillas

Strings sin cerrar

Faltas de comas

Uso incorrecto de : o {, [

📌 Conclusiones
El proyecto ha permitido reforzar conceptos clave como análisis sintáctico, recursión, y manejo de errores estructurados.

Go es un lenguaje robusto para trabajar con análisis de texto, gracias a su control de errores explícito y tipos seguros.

Desarrollar un parser desde cero ayuda a comprender cómo funcionan internamente bibliotecas como encoding/json.

🔧 Autores: Jhonatan Jair Huaman Yovera y Claudia Regalado Diaz
🎓 Curso: Taller de Lenguajes de Programación
📅 Fecha: Julio 2025
