package ai

import (
	"fmt"
	lib "kxmj.common/mahjong"
	"testing"
)

func TestShanten_CalculateMinShanten(t *testing.T) {
	type fields struct {
		handLen         int32
		handCards       []int32
		shunKeNums      int32
		daZiNums        int32
		pairNums        int32
		mustOutCardNums int32
		signalCard      int64
		anGangCard      int64
		minShanten      int32
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shanten{
				handLen:         tt.fields.handLen,
				handCards:       tt.fields.handCards,
				shunKeNums:      tt.fields.shunKeNums,
				daZiNums:        tt.fields.daZiNums,
				pairNums:        tt.fields.pairNums,
				mustOutCardNums: tt.fields.mustOutCardNums,
				signalCard:      tt.fields.signalCard,
				anGangCard:      tt.fields.anGangCard,
				minShanten:      tt.fields.minShanten,
			}
			if got := s.CalculateShanten(); got != tt.want {
				t.Errorf("CalculateMinShanten() = %v, want %v", got, tt.want)
			}
		})
	}
}

var AllCards = lib.Cards{
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // 1-9万
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // 1-9万
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // 1-9万
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // 1-9万

	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, // 1-9条
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, // 1-9条
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, // 1-9条
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, // 1-9条

	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, // 1-9筒
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, // 1-9筒
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, // 1-9筒
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, // 1-9筒

	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, //东南西北中发白
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, //东南西北中发白
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, //东南西北中发白
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, //东南西北中发白
}

func BenchmarkTest(b *testing.B) {
	stack := lib.NewCardStack(AllCards)
	for i := 0; i < b.N; i++ {
		handCards, _ := stack.PickUpCards(14)
		handCards.Sort()
		fmt.Println(handCards)
		shanten := NewShanten(CardToTiles(handCards))
		//fmt.Println(shanten)
		shanten.DecomposeHandTile(0)
		//if shanten.minShanten > 1 {
		//	handCards.Sort()
		//	fmt.Printf("handCards:%v ,shanten:%v\n", handCards, shanten)
		//}
		if stack.GetResidueCardsNum() < 14 {
			stack.Reset()
		}
	}
}

func TestShanten_DecomposeDNXB(t *testing.T) {
	type fields struct {
		handLen         int32
		handCards       []int32
		shunKeNums      int32
		daZiNums        int32
		pairNums        int32
		mustOutCardNums int32
		signalCard      int64
		anGangCard      int64
		minShanten      int32
	}
	type args struct {
		depth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shanten{
				handLen:         tt.fields.handLen,
				handCards:       tt.fields.handCards,
				shunKeNums:      tt.fields.shunKeNums,
				daZiNums:        tt.fields.daZiNums,
				pairNums:        tt.fields.pairNums,
				mustOutCardNums: tt.fields.mustOutCardNums,
				signalCard:      tt.fields.signalCard,
				anGangCard:      tt.fields.anGangCard,
				minShanten:      tt.fields.minShanten,
			}
			s.DecomposeDNXB(tt.args.depth)
		})
	}
}

func TestShanten_DecomposeHandTile(t *testing.T) {
	type fields struct {
		handLen         int32
		handCards       []int32
		shunKeNums      int32
		daZiNums        int32
		pairNums        int32
		mustOutCardNums int32
		signalCard      int64
		anGangCard      int64
		minShanten      int32
	}
	type args struct {
		depth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shanten{
				handLen:         tt.fields.handLen,
				handCards:       tt.fields.handCards,
				shunKeNums:      tt.fields.shunKeNums,
				daZiNums:        tt.fields.daZiNums,
				pairNums:        tt.fields.pairNums,
				mustOutCardNums: tt.fields.mustOutCardNums,
				signalCard:      tt.fields.signalCard,
				anGangCard:      tt.fields.anGangCard,
				minShanten:      tt.fields.minShanten,
			}
			s.DecomposeHandTile(tt.args.depth)
		})
	}
}

func TestShanten_DecomposeWTT(t *testing.T) {
	type fields struct {
		handLen         int32
		handCards       []int32
		shunKeNums      int32
		daZiNums        int32
		pairNums        int32
		mustOutCardNums int32
		signalCard      int64
		anGangCard      int64
		minShanten      int32
	}
	type args struct {
		depth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shanten{
				handLen:         tt.fields.handLen,
				handCards:       tt.fields.handCards,
				shunKeNums:      tt.fields.shunKeNums,
				daZiNums:        tt.fields.daZiNums,
				pairNums:        tt.fields.pairNums,
				mustOutCardNums: tt.fields.mustOutCardNums,
				signalCard:      tt.fields.signalCard,
				anGangCard:      tt.fields.anGangCard,
				minShanten:      tt.fields.minShanten,
			}
			s.DecomposeWTT(tt.args.depth)
		})
	}
}

func TestShanten_DecomposeZFB(t *testing.T) {
	type fields struct {
		handLen         int32
		handCards       []int32
		shunKeNums      int32
		daZiNums        int32
		pairNums        int32
		mustOutCardNums int32
		signalCard      int64
		anGangCard      int64
		minShanten      int32
	}
	type args struct {
		depth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shanten{
				handLen:         tt.fields.handLen,
				handCards:       tt.fields.handCards,
				shunKeNums:      tt.fields.shunKeNums,
				daZiNums:        tt.fields.daZiNums,
				pairNums:        tt.fields.pairNums,
				mustOutCardNums: tt.fields.mustOutCardNums,
				signalCard:      tt.fields.signalCard,
				anGangCard:      tt.fields.anGangCard,
				minShanten:      tt.fields.minShanten,
			}
			s.DecomposeZFB(tt.args.depth)
		})
	}
}

func TestShanten_String(t *testing.T) {
	type fields struct {
		handLen         int32
		handCards       []int32
		shunKeNums      int32
		daZiNums        int32
		pairNums        int32
		mustOutCardNums int32
		signalCard      int64
		anGangCard      int64
		minShanten      int32
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shanten{
				handLen:         tt.fields.handLen,
				handCards:       tt.fields.handCards,
				shunKeNums:      tt.fields.shunKeNums,
				daZiNums:        tt.fields.daZiNums,
				pairNums:        tt.fields.pairNums,
				mustOutCardNums: tt.fields.mustOutCardNums,
				signalCard:      tt.fields.signalCard,
				anGangCard:      tt.fields.anGangCard,
				minShanten:      tt.fields.minShanten,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
