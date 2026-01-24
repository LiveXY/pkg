package reflectx

import (
	"testing"
)

type User struct {
	ID   int    `gorm:"primaryKey;column:id"`
	Name string `gorm:"column:name"`
	Age  int
}

type Order struct {
	OrderID int `gorm:"primaryKey"`
	Price   float64
}

// TestGetTableStruct 测试通过反射提取结构体元数据（如表名、字段数、主键标记）的功能
func TestGetTableStruct(t *testing.T) {
	// Test User
	table := GetTableStruct[User]()
	if table.Name != "User" {
		t.Errorf("期望名称为 User，实际结果为 %v", table.Name)
	}
	if table.Fields != 3 {
		t.Errorf("期望字段数为 3，实际结果为 %d", table.Fields)
	}
	if len(table.PKs) != 1 || table.PKs[0] != "id" {
		t.Errorf("期望主键为 id，实际结果为 %v", table.PKs)
	}

	// Test Order
	tableOrder := GetTableStruct[Order]()
	if len(tableOrder.PKs) != 1 || tableOrder.PKs[0] != "OrderID" {
		t.Errorf("期望主键为 OrderID，实际结果为 %v", tableOrder.PKs)
	}
}
