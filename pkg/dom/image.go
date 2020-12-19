package dom

type Image struct {
	elem *Elem
}

func AsImage(elem *Elem) *Image {
	if elem.Tag() != "img" {
		return nil
	}
	return &Image{
		elem: elem,
	}
}

func (img *Image) SetSrc(src string) {
	img.elem.obj.Set("src", img.elem.factory.String(src))
}
