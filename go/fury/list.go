package fury

import "reflect"

type listSerializer struct {
}

const (
	// If track elements ref, use the first bit 0b1 of the header to flag it.
	TrackingRefFlag int8 = 0b1
	// If the elements have null, use the second bit 0b10 of the header to flag it. If ref tracking is enabled for this element type, this flag is invalid.
	HaveNull int8 = 0b10
	// If the element types are not the declared type, use the 3rd bit 0b100 of the header to flag it.
	NotDeclaredType int8 = 0b100
	// If the element types are different, use the 4rd bit 0b1000 header to flag it.
	NotTheSameType int8 = 0b1000
)

func (s listSerializer) TypeId() TypeId {
	return LIST
}

type header struct {
	isHomogeneous  bool
	isNotNull      bool
	isDeclaredType bool
}

func (s listSerializer) checkHeader(value reflect.Value) *header {
	length := value.Len()

	result := header{
		isHomogeneous:  true,
		isNotNull:      true,
		isDeclaredType: true,
	}

	var firstElementType = value.Index(0).Type()
	for i := 1; i < length; i++ {
		nValue := value.Index(i)
		if nValue.IsNil() {
			result.isNotNull = false
		}
		if nValue.Type() != firstElementType {
			result.isHomogeneous = false
		}

		if !result.isNotNull && !result.isHomogeneous {
			break
		}
	}
	return &result
}

func (s listSerializer) Write(f *Fury, buf *ByteBuffer, value reflect.Value) error {
	length := value.Len()

	if err := f.writeLength(buf, length); err != nil {
		return err
	}

	header := s.checkHeader(value)

	flag := int8(0)

	if f.referenceTracking {
		flag |= TrackingRefFlag
	}
	if !header.isHomogeneous {
		flag |= NotTheSameType
	}
	if !header.isNotNull {
		flag |= HaveNull
	}
	if !header.isDeclaredType {
		flag |= NotDeclaredType
	}

	if err := buf.WriteByte(byte(flag)); err != nil {
		return err
	}

    if (header.isHomogeneous) {
        vType := value.Index(0).Type()
    }

	for i := 0; i < length; i++ {
		if err := f.WriteReferencable(buf, value.Index(i)); err != nil {
			return err
		}
	}
	return nil
}
func (s listSerializer) Read(f *Fury, buf *ByteBuffer, type_ reflect.Type, value reflect.Value) error {
	length := f.readLength(buf)
	if value.Cap() < length {
		value.Set(reflect.MakeSlice(value.Type(), length, length))
	} else if value.Len() < length {
		value.Set(value.Slice(0, length))
	}
	f.refResolver.Reference(value)
	for i := 0; i < length; i++ {
		elem := value.Index(i)
		if err := f.ReadReferencable(buf, elem); err != nil {
			return err
		}
	}
	return nil
}
