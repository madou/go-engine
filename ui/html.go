package ui

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/aymerick/douceur/css"
	"github.com/aymerick/douceur/parser"
	vmath "github.com/walesey/go-engine/vectormath"
	"golang.org/x/net/html"
)

// LoadPage - load the html/css code into the container
func LoadPage(container *Container, htmlInput, cssInput io.Reader, assets HtmlAssets) error {
	document, err := html.Parse(htmlInput)
	if err != nil {
		log.Printf("Error parsing html: %v", err)
		return err
	}

	css, err := ioutil.ReadAll(cssInput)
	if err != nil {
		log.Printf("Error reading css: %v", err)
		return err
	}

	styles, err := parser.Parse(string(css))
	if err != nil {
		log.Printf("Error parsing css: %v", err)
		return err
	}

	renderNode(container, document.FirstChild, styles, assets)

	return nil
}

func renderNode(container *Container, node *html.Node, styles *css.Stylesheet, assets HtmlAssets) {
	nextNode := node
	for nextNode != nil {
		if nextNode.Type == 1 {
			text := nextNode.Data
			text = strings.TrimSpace(text)
			if len(text) > 0 {
				textElement := NewTextElement(text, color.Black, 32, assets.fontMap["default"])
				container.AddChildren(textElement)
			}
		} else {
			newContainer := NewContainer()
			container.AddChildren(newContainer)
			//Parse Styles
			normalStyles := getStyles(styles, nextNode, "")
			applyStyles(newContainer, normalStyles)
			hoverStyles := getStyles(styles, nextNode, ":hover")
			if len(hoverStyles) > 0 {
				newContainer.Hitbox.AddOnHover(func() {
					applyDefaultStyles(newContainer)
					applyStyles(newContainer, normalStyles)
					applyStyles(newContainer, hoverStyles)
				})
				newContainer.Hitbox.AddOnUnHover(func() {
					applyDefaultStyles(newContainer)
					applyStyles(newContainer, normalStyles)
				})
			}
			//Parse Attributes
			for _, attr := range nextNode.Attr {
				switch {
				case attr.Key == "onclick":
					callback, ok := assets.callbackMap[attr.Val]
					if ok {
						newContainer.Hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
							callback(newContainer, button, release, position)
						})
					}
				case attr.Key == "onhover":
					callback, ok := assets.callbackMap[attr.Val]
					if ok {
						newContainer.Hitbox.AddOnHover(func() {
							callback(newContainer)
						})
					}
				case attr.Key == "onunhover":
					callback, ok := assets.callbackMap[attr.Val]
					if ok {
						newContainer.Hitbox.AddOnUnHover(func() {
							callback(newContainer)
						})
					}
				}
			}
			//tag type
			if nextNode.DataAtom.String() == "input" {
				inputType := getAttribute(nextNode, "type")
				switch {
				case inputType == "text":
					newContainer.AddChildren(NewTextElement("", color.Black, 32, assets.fontMap["default"]))
				}
			}
			renderNode(newContainer, nextNode.FirstChild, styles, assets)
		}
		if nextNode == nextNode.NextSibling {
			break
		}
		nextNode = nextNode.NextSibling
	}
	// fmt.Println(styles.Rules[0].Selectors)
}

func applyDefaultStyles(container *Container) {
	container.SetBackgroundColor(0, 0, 0, 0)
	container.SetHeight(0)
	container.SetWidth(0)
	container.SetMargin(NewMargin(0))
	container.SetPadding(NewMargin(0))
}

func applyStyles(container *Container, styles map[string]string) {
	for prop, value := range styles {
		switch {
		case prop == "padding":
			paddings := parseDimensions(value)
			if len(paddings) == 1 {
				container.SetPadding(NewMargin(paddings[0]))
			} else if len(paddings) == 4 {
				container.SetPadding(Margin{paddings[0], paddings[1], paddings[2], paddings[3]})
			}
		case prop == "margin":
			margins := parseDimensions(value)
			if len(margins) == 1 {
				container.SetMargin(NewMargin(margins[0]))
			} else if len(margins) == 4 {
				container.SetMargin(Margin{margins[0], margins[1], margins[2], margins[3]})
			}
		case prop == "background-color":
			color := parseColor(value)
			container.SetBackgroundColor(color[0], color[1], color[2], color[3])
		case prop == "width":
			width := parseDimensions(value)
			if len(width) == 1 {
				container.SetWidth(width[0])
			}
		case prop == "height":
			height := parseDimensions(value)
			if len(height) == 1 {
				container.SetHeight(height[0])
			}
		default:
			fmt.Printf("Unsupported css prop: %v: %v;\n", prop, value)
		}
	}
}

func parseDimensions(dimensionsStr string) []float64 {
	dimensions := strings.Split(dimensionsStr, " ")
	values := make([]float64, len(dimensions))
	for i, dimension := range dimensions {
		value, err := strconv.ParseFloat(strings.Replace(dimension, "px", "", 1), 64) // TODO fix this
		if err != nil {
			fmt.Printf("Error parsing dimensions: %v;\n", err)
		}
		values[i] = value
	}
	return values
}

func parseColor(colorStr string) [4]uint8 {
	data := []byte(colorStr)
	r, g, b, a := "0", "0", "0", "ff"
	if len(data) == 4 {
		r, g, b = string(data[1:2]), string(data[2:3]), string(data[3:])
		r, g, b = fmt.Sprintf("%v0", r), fmt.Sprintf("%v0", g), fmt.Sprintf("%v0", b)
	} else if len(data) == 5 {
		r, g, b, a = string(data[1:2]), string(data[2:3]), string(data[3:4]), string(data[4:])
		r, g, b, a = fmt.Sprintf("%v0", r), fmt.Sprintf("%v0", g), fmt.Sprintf("%v0", b), fmt.Sprintf("%v0", a)
	} else if len(data) == 7 {
		r, g, b = string(data[1:3]), string(data[3:5]), string(data[5:])
	} else if len(data) == 9 {
		r, g, b, a = string(data[1:3]), string(data[3:5]), string(data[5:7]), string(data[7:])
	}
	rd, _ := hex.DecodeString(r)
	gd, _ := hex.DecodeString(g)
	bd, _ := hex.DecodeString(b)
	ad, _ := hex.DecodeString(a)
	result := [4]uint8{0, 0, 0, 0}
	if len(rd) > 0 {
		result[0] = uint8(rd[0])
	}
	if len(gd) > 0 {
		result[1] = uint8(gd[0])
	}
	if len(bd) > 0 {
		result[2] = uint8(bd[0])
	}
	if len(ad) > 0 {
		result[3] = uint8(ad[0])
	}
	return result
}

func getStyles(styles *css.Stylesheet, node *html.Node, modifier string) map[string]string {
	id := getAttribute(node, "id")
	class := getAttribute(node, "class")
	tag := node.DataAtom.String()
	rules := make(map[string]string)
	getStyleBySelector(styles, rules, fmt.Sprintf("#%v%v", id, modifier))
	getStyleBySelector(styles, rules, fmt.Sprintf(".%v%v", class, modifier))
	getStyleBySelector(styles, rules, fmt.Sprintf("%v%v", tag, modifier))
	// TODO: get syles from node.Parent
	return rules
}

func getAttribute(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getStyleBySelector(styles *css.Stylesheet, dest map[string]string, selector string) {
	for _, rule := range styles.Rules {
		for _, sel := range rule.Selectors {
			if sel == selector {
				for _, declaration := range rule.Declarations {
					if _, ok := dest[declaration.Property]; !ok {
						dest[declaration.Property] = declaration.Value
					}
				}
			}
		}
	}
}
