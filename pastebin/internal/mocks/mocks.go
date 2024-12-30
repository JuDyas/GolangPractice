package mocks

//go:generate mockgen -source=../servises/paste.go -destination=./service_mocks.go -package=mocks github.com/JuDyas/GolangPractice/pastebin/internal/mocks Service
