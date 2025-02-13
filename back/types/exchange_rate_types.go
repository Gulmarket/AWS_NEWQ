package types

type RSS struct {
    Items []Item `xml:"item"`
}

type Item struct {
    FullName    string `xml:"fullname"`
    Title       string `xml:"title"`
    Description string `xml:"description"`
    Quant       string `xml:"quant"`
    Index       string `xml:"index"`
    Change      string `xml:"change"`
}
