package storage

import (
	"github.com/ad-sho-loko/bogodb/meta"
	"github.com/gogo/protobuf/proto"
)

func NewTuple(minTxId uint64, values []interface{}) *Tuple{
	var t Tuple
	t.MinTxId = minTxId
	t.MaxTxId = minTxId

	var td *TupleData
	for _, v := range values{
		switch concrete := v.(type) {

		case int:
			td = &TupleData{
				Type:TupleData_INT,
				Number:*proto.Int32(int32(concrete)),
			}

		case string:
			td = &TupleData{
				Type:TupleData_STRING,
				String_:*proto.String(concrete),
			}
		}

		t.Data = append(t.Data, td)
	}

	return &t
}

func (m *Tuple) Less(than meta.Item) bool{
	t, ok := than.(*Tuple)

	if !ok{
		return false
	}

	return m.GetPk() < t.GetPk()
}

func (m *Tuple) GetPk() int{
	return 0
}

func SerializeTuple(t *Tuple) ([128]byte, error){
	out, err := proto.Marshal(t)

	if err != nil{
		return [128]byte{}, err
	}

	var b [128]byte
	copy(b[:], out)

	return b, nil
}

func DeserializeTuple(b [128]byte) (*Tuple, error){
	var t Tuple

	err := proto.Unmarshal(b[:], &t)
	if err != nil{
		// return nil, err
	}

	return &t, nil
}

func (m *Tuple) GetInt(table *meta.Table, colName string) int{
	for _, col := range table.Columns{
		if colName == col.Name{
		}
	}
	return 0
}

func (m *Tuple) GetStr(table string, col string) string{
	return ""
}

func (m *Tuple) IsUnused() bool{
	// If minTxId is zero, it's an empty tuple.
	return m.MinTxId == 0
}

func (m *Tuple) CanSee(tran *Transaction) bool{
	if m.MinTxId == tran.txid{
		return true
	}

	if m.MaxTxId < tran.Txid(){
		return false
	}

	return true
}
