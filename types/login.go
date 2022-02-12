package types


type LoginDTO struct {
	Login string;
	Password string;
}

type LoginDB struct {
	Login string;
	Hash string;
	Role uint8
}

type RegisterDTO struct {
	LoginDTO
	Role uint8
}