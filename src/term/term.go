package term

type Term struct {
    key string
    value string
    start int
    end int
}

type TermFrequency struct {
    Term
    frequency int
}

func (t *Term) GetKey() string {
    return t.key
}

func (t *Term) GetValue() string {
    return t.value
}

//set the value as newValue and return old value
func (t *Term) SetValue(newValue string) string {
    oldValue := t.value
    t.value = newValue
    return oldValue
}

func (t *Term) GetStart() int {
    return t.start
}

func (t *Term) GetEnd() int {
    return t.end
}

func (tf *TermFrequency) GetFrequency() int {
    return tf.frequency
}

func (tf *TermFrequency) Increase() {
    tf.frequency++
}

func NewTerm(key, value, start, end) *Term {
    return &Term {
        key: key,
        value: value,
        start: start,
        end: end
    }
}
