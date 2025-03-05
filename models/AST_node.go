package models

type ASTNode struct {
    Operator      string    // Оператор (+, -, *, /), если это не число
    Value         float64   // Значение, если это число
    IsLeaf        bool      // true, если это число (лист)
    TaskScheduled bool      // true, если задача для этого узла уже создана
    Left          *ASTNode  // Ссылка на левый узел (может быть nil)
    Right         *ASTNode  // Ссылка на правый узел (может быть nil)
}
