package usersUsecases

import (
	"github.com/SurrealGod/jelly-shop-code/config"
	"github.com/SurrealGod/jelly-shop-code/modules/users"
	"github.com/SurrealGod/jelly-shop-code/modules/users/usersRepositories"
	usersrepositories "github.com/SurrealGod/jelly-shop-code/modules/users/usersRepositories"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersrepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
  if err :=req.BcryptHashing(); err != nil {
	 return nil,err
  }


  //insert user
  result,err := u.usersRepository.InsertUser(req,false)
  if err != nil {
	return nil,err
  }
	
	return result,nil
}
