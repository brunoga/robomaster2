package internal

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/brunoga/robomaster2/internal/callbacks"
	"github.com/brunoga/robomaster2/internal/dji"
	"github.com/brunoga/robomaster2/internal/event"
	"github.com/brunoga/robomaster2/internal/unitybridge"
)

type ResultHandler func(result *dji.Result, wg *sync.WaitGroup)

type CommandController struct {
	*GenericController

	startListeningCallbacks *callbacks.Callbacks
	performActionCallbacks  *callbacks.Callbacks
	setValueForKeyCallbacks *callbacks.Callbacks
}

func NewCommandController() (*CommandController, error) {
	startListeningCallbacks := callbacks.New(
		"CommandController/StartListening",
		func(key callbacks.Key) error {
			ev := event.NewWithSubType(event.TypeStartListening, uint64(key))

			unitybridge.Get().SendEvent(int64(ev.Code()), nil, 0)

			return nil
		},
		func(key callbacks.Key) error {
			ev := event.NewWithSubType(event.TypeStopListening, uint64(key))

			unitybridge.Get().SendEvent(int64(ev.Code()), nil, 0)

			return nil
		},
	)

	performActionCallbacks := callbacks.New(
		"CommandController/PerformAction", nil, nil)

	setValueForKeyCallbacks := callbacks.New(
		"CommandController/SetValueForKey", nil, nil)

	cc := &CommandController{
		startListeningCallbacks: startListeningCallbacks,
		performActionCallbacks:  performActionCallbacks,
		setValueForKeyCallbacks: setValueForKeyCallbacks,
	}

	cc.GenericController = NewGenericController(cc)

	var err error
	err = cc.StartControllingEvent(event.TypeGetValue)
	if err != nil {
		return nil, err
	}
	err = cc.StartControllingEvent(event.TypeSetValue)
	if err != nil {
		return nil, err
	}
	err = cc.StartControllingEvent(event.TypePerformAction)
	if err != nil {
		return nil, err
	}
	err = cc.StartControllingEvent(event.TypeStartListening)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func (c *CommandController) StartListening(key dji.Key,
	resultHandler ResultHandler) (uint64, error) {
	if key < 1 || key >= dji.KeysCount {
		return 0, fmt.Errorf("invalid key")
	}
	if resultHandler == nil {
		return 0, fmt.Errorf("eventHandler must not be nil")
	}
	if (key.AccessType() & dji.KeyAccessTypeRead) == 0 {
		return 0, fmt.Errorf("key is not readable")
	}

	tag, err := c.startListeningCallbacks.AddContinuous(callbacks.Key(
		key.Value()), resultHandler)

	return uint64(tag), err
}

func (c *CommandController) StopListening(key dji.Key, tag uint64) error {
	if key < 1 || key >= dji.KeysCount {
		return fmt.Errorf("invalid key")
	}

	return c.startListeningCallbacks.Remove(callbacks.Key(key),
		callbacks.Tag(tag))
}

func (c *CommandController) PerformAction(key dji.Key, param interface{},
	resultHandler ResultHandler) error {
	if key < 1 || key >= dji.KeysCount {
		return fmt.Errorf("invalid key")
	}
	if (key.AccessType() & dji.KeyAccessTypeAction) == 0 {
		return fmt.Errorf("key can not be acted upon")
	}

	var err error
	tag := callbacks.Tag(0)

	if resultHandler != nil {
		tag, err = c.performActionCallbacks.AddSingleShot(
			callbacks.Key(key.Value()), resultHandler)
		if err != nil {
			return err
		}
	}

	var data []byte
	if param != nil {
		data, err = json.Marshal(param)
		if err != nil {
			return err
		}
	}

	ev := event.NewWithSubType(event.TypePerformAction, uint64(key.Value()))
	unitybridge.Get().SendEvent(int64(ev.Code()), data, int64(tag))

	return nil
}

func (c *CommandController) SetValueForKey(key dji.Key, param interface{},
	resultHandler ResultHandler) error {
	if key < 1 || key >= dji.KeysCount {
		return fmt.Errorf("invalid key")
	}
	if (key.AccessType() & dji.KeyAccessTypeWrite) == 0 {
		return fmt.Errorf("key can not be written to")
	}

	var err error
	tag := callbacks.Tag(0)

	if resultHandler != nil {
		tag, err = c.setValueForKeyCallbacks.AddSingleShot(
			callbacks.Key(key.Value()), resultHandler)
		if err != nil {
			return err
		}
	}

	var data []byte
	if param != nil {
		data, err = json.Marshal(param)
		if err != nil {
			return err
		}
	}

	ev := event.NewWithSubType(event.TypeSetValue, uint64(key.Value()))

	unitybridge.Get().SendEventWithString(int64(ev.Code()), string(data), int64(tag))

	return nil
}

func (c *CommandController) DirectSendValue(key dji.Key, value uint64) {
	ev := event.NewWithSubType(event.TypePerformAction, uint64(key.Value()))
	unitybridge.Get().SendEventWithNumber(int64(ev.Code()), int64(value), 0)
}

func (c *CommandController) HandleEventCallback(code int64, data []byte,
	tag int64) {
	var value interface{}

	// TODO(bga): Apparently, the unity bridge reserves the upper 8 bits
	//  for reporting back type information. Double check this.
	dataType := tag >> 56
	switch dataType {
	case 0:
		value = string(data)
	case 1:
		value = binary.LittleEndian.Uint64(data)
	default:
		// Apparently only string and uint64 types are supported
		// currently.
		panic(fmt.Sprintf("Unexpected data type: %d.\n", dataType))
	}

	// See above.
	adjustedTag := tag & 0xffffffffffffff

	ev := event.NewEventFromCode(uint64(code))

	switch ev.Type() {
	case event.TypeStartListening:
		c.handleStartListening(ev.SubType(), value)
	case event.TypePerformAction:
		c.handlePerformAction(ev.SubType(), value, uint64(adjustedTag))
	case event.TypeSetValue:
		c.handleSetValue(ev.SubType(), value, uint64(adjustedTag))
	default:
		log.Printf("Unsupported event %s. Value:%v. Tag:%d\n",
			event.TypeName(ev.Type()), value, tag)
	}
}

func (c *CommandController) handleStartListening(key uint64,
	value interface{}) {
	stringValue, ok := value.(string)
	if !ok {
		panic("unexpected non-string value")
	}

	cbs, err := c.startListeningCallbacks.CallbacksForKey(
		callbacks.Key(key))
	if err != nil {
		log.Printf("Error looking up callbacks for key %d: %s\n", key,
			err)
		return
	}

	result := dji.NewResultFromJSON([]byte(stringValue))

	var wg sync.WaitGroup
	for _, cb := range cbs {
		wg.Add(1)
		go cb.(ResultHandler)(result, &wg)
	}

	wg.Wait()
}

func (c *CommandController) handlePerformAction(key uint64,
	value interface{}, tag uint64) {
	stringValue, ok := value.(string)
	if !ok {
		panic("unexpected non-string value")
	}

	cb, err := c.startListeningCallbacks.Callback(callbacks.Key(key),
		callbacks.Tag(tag))
	if err != nil {
		// No callback set. That is fine.
		return
	}

	result := dji.NewResultFromJSON([]byte(stringValue))

	cb.(ResultHandler)(result, nil)
}

func (c *CommandController) handleSetValue(key uint64,
	value interface{}, tag uint64) {
	stringValue, ok := value.(string)
	if !ok {
		panic("unexpected non-string value")
	}

	cb, err := c.startListeningCallbacks.Callback(callbacks.Key(key),
		callbacks.Tag(tag))
	if err != nil {
		// No callback set. That is fine.
		return
	}

	result := dji.NewResultFromJSON([]byte(stringValue))

	cb.(ResultHandler)(result, nil)
}

func (c *CommandController) Teardown() error {
	var err error
	err = c.StopControllingEvent(event.TypeGetValue)
	if err != nil {
		return err
	}
	err = c.StopControllingEvent(event.TypeSetValue)
	if err != nil {
		return err
	}
	err = c.StopControllingEvent(event.TypePerformAction)
	if err != nil {
		return err
	}
	err = c.StopControllingEvent(event.TypeStartListening)
	if err != nil {
		return err
	}

	// TODO(bga): Also clean up ResultHandler callbacks before returning.

	return nil
}
