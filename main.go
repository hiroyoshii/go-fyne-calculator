package main

import (
    "fmt"
    "os"
    "bufio"
    "io/ioutil"
    "strings"
    "go/token"
    "go/types"

    "fyne.io/fyne"
    "fyne.io/fyne/app"
    "fyne.io/fyne/widget"
    "fyne.io/fyne/layout"
)

type calculator struct {
  displayText *widget.Label
  equation string
  operated bool
}

func (c *calculator) operator(op string){
  c.equation = c.equation + op
  c.operated = true
}

func (c *calculator) evaluate(){
  if c.operated {
    return
  }
  eval, err := types.Eval(token.NewFileSet(), nil, token.NoPos, c.equation)
  if err != nil {
      fmt.Println(err)
      return
  }
  c.displayText.SetText(fmt.Sprintf("%s",eval.Value))
  c.equation = ""
  c.operated = true
}

func (c *calculator) clearAll(){
  c.equation = ""
  c.displayText.SetText(c.equation)
  c.operated = false
}
func (c *calculator) backspace(){
  if c.operated || len(c.equation)-1 < 0{
    return
  }
  c.equation = string([]rune(c.equation)[:len(c.equation)-1])
  c.displayText.SetText(c.equation)
}

func (c *calculator) number(n string){
  c.equation = c.equation + n
  if c.operated {
    c.displayText.SetText(n)
  } else {
    c.displayText.SetText(c.equation)
  }
  c.operated = false
}

func (c *calculator) point() {
  if c.operated {
    return
  }
  c.equation += string(".")
  c.displayText.SetText(c.equation)
}

func (c *calculator) togglePlusMinus() {
  if c.operated {
    return
  }
  if strings.Contains(c.equation, "-") {
    c.equation =  string([]rune(c.equation)[1:])
    c.displayText.SetText(c.equation)
  } else {
    c.equation = string("-") + c.equation
    c.displayText.SetText(c.equation)
  }
}

func main() {
    a := app.New()
    w := a.NewWindow("Timer")

    calc := calculator{}
    calc.displayText = &widget.Label{Text:""}
    calc.displayText.Alignment = fyne.TextAlignTrailing
    calc.displayText.TextStyle.Monospace = true

    clearGrid := fyne.NewContainerWithLayout(layout.NewGridLayout(1), newClearGrid(&calc), newNumGrid(&calc))
    buttons := fyne.NewContainerWithLayout(layout.NewGridLayout(2), clearGrid, newOperatorGrid(&calc))
    container := fyne.NewContainerWithLayout(layout.NewGridLayout(1), calc.displayText, buttons)
    w.SetContent(container)
    w.ShowAndRun()
}

func newNumGrid(calc *calculator) *fyne.Container {
  var blocks []fyne.CanvasObject

  for i:=1; i<=3*3; i++{
    strNum := fmt.Sprintf("%d", i)
    blocks = append(blocks, &widget.Button{
        Text: strNum, OnTapped: func() { calc.number(strNum) },
    })
  }
  blocks = append(blocks, &widget.Button{
      Text: "+/-", OnTapped: func() { calc.togglePlusMinus() },
  })
  blocks = append(blocks, &widget.Button{
      Text: "0", OnTapped: func() { calc.number("0") },
  })
  blocks = append(blocks, &widget.Button{
      Text: ".", OnTapped: func() { calc.point() },
  })

  return fyne.NewContainerWithLayout(layout.NewGridLayout(3), blocks...)
}

func newOperatorGrid(calc *calculator) *fyne.Container {
  var operators []fyne.CanvasObject

  operators = append(operators, &widget.Button{
      Text: "/", OnTapped: func() { calc.operator("/") },
  })
  operators = append(operators, &widget.Button{
      Text: "*", OnTapped: func() { calc.operator("*") },
  })
  operators = append(operators, &widget.Button{
      Text: "-", OnTapped: func() { calc.operator("-") },
  })
  operators = append(operators, &widget.Button{
      Text: "+", OnTapped: func() { calc.operator("+") },
  })
  operators = append(operators, &widget.Button{
      Text: "=", OnTapped: func() { calc.evaluate() },
  })

  return fyne.NewContainerWithLayout(layout.NewGridLayout(1), operators...)
}
func newClearGrid(calc *calculator) *fyne.Container {
  var clears []fyne.CanvasObject

  clears = append(clears, &widget.Button{
      Text: "C", OnTapped: func() { calc.clearAll() },
  })

  iconFile,_ := os.Open("images/backspace.png")
  r := bufio.NewReader(iconFile)
  b,_ := ioutil.ReadAll(r)
  clears = append(clears, &widget.Button{
      Icon: fyne.NewStaticResource("icon", b), OnTapped: func() { calc.backspace() },
  })

  return fyne.NewContainerWithLayout(layout.NewGridLayout(2), clears...)
}
