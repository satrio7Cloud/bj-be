package utils

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func HashPassword(planiPassword string) (string, error) {
	hash, err := argon2id.CreateHash(planiPassword, argon2id.DefaultParams)
	if err != nil {
		log.Printf("Error saat hashing password: %v", err)
		return "", err
	}

	return hash, nil
}

func ComparePassword(planPassword, hashPassword string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(planPassword, hashPassword)
	if err != nil {
		log.Printf("Error saat membandingkan password: %v", err)
		return false, err
	}
	return match, nil

}
