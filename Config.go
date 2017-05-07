package main

type Config struct {
	SpriteName string `yaml:"spriteName"`
	Images []Image `yaml:"images"`
}

type Image struct {
	FileName string `yaml:"fileName"`
}