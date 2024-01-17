package models

type ApprovalUsecase interface {
	Update(int,*UpdateStatusReq) (*Approval, error)
	ReceiveRequest(id int,optional map[string]interface{}) ([]Approval, error)
}

type ApprovalRepository interface {
}

type ProducerProfile interface {
}
