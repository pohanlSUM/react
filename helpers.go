// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package react

import (
	"github.com/gopherjs/gopherjs/js"
)

// Fragment is used to group a list of children
// without adding extra nodes to the DOM.
//
// See: https://reactjs.org/docs/fragments.html
func Fragment(key *string, children ...interface{}) *js.Object {
	props := map[string]interface{}{}
	if key != nil {
		props["key"] = *key
	}
	return JSX(React.Get("Fragment"), props, children...)
}

// OnRenderCallback is the callback function signature of the onRender argument to Profiler function.
type OnRenderCallback func(id string, phase string, actualDuration, baseDuration float64, startTime, commitTime float64, interactions *js.Object)

// Profiler is used to find performance bottlenecks in your application.
//
// See: https://reactjs.org/docs/profiler.html
func Profiler(id string, onRender OnRenderCallback, children ...interface{}) *js.Object {

	props := map[string]interface{}{
		"id": id,
	}
	if onRender != nil {
		props["onRender"] = onRender
	}

	return JSX(React.Get("Profiler"), props, children...)
}

// JSX is used to create an Element.
func JSX(component interface{}, props interface{}, children ...interface{}) *js.Object {

	args := []interface{}{
		component,
		SToMap(props),
	}
	if len(children) > 0 {
		args = append(args, children...)
	}

	return React.Call("createElement", args...)
}

// JSFn is a convenience function used to call javascript native functions.
// If the native function throws an exception, then a *js.Error is returned.
//
// Example:
//
//  // alert('Hello World!')
//  JSFn("alert", "Hello World!")
//
//  // JSON.parse('{"name":"John"}')
//  JSFn("JSON.parse", `{"name":"John"}`)
//
func JSFn(funcName string, args ...interface{}) (_ *js.Object, rErr error) {
	defer func() {
		if e := recover(); e != nil {
			err, ok := e.(*js.Error)
			if !ok {
				panic(e)
			}
			rErr = err
		}
	}()

	out := js.Global

	splits := js.Global.Get("String").Invoke(funcName).Call("split", ".")
	for idx := 0; idx < splits.Length(); idx++ {
		split := splits.Index(idx).String()
		if idx == splits.Length()-1 {
			out = out.Call(split, args...)
		} else {
			out = out.Get(split)
		}
	}

	return out, nil
}

// CreateRef will create a Ref.
//
// See: https://reactjs.org/docs/refs-and-the-dom.html
func CreateRef() *js.Object {
	return React.Call("createRef")
}

// ForwardRef will forward a Ref to child components.
//
// See: https://reactjs.org/docs/forwarding-refs.html
func ForwardRef(component interface{}) *js.Object {
	return React.Call("forwardRef", func(props *js.Object, ref *js.Object) *js.Object {
		props.Set("ref", ref)

		n := React.Get("Children").Call("count", props.Get("children")).Int()
		switch n {
		case 0:
			return JSX(component, props)
		case 1:
			return JSX(component, props, props.Get("children"))
		default:
			children := []interface{}{}
			for i := 0; i < n; i++ {
				children = append(children, props.Get("children").Index(i))
			}
			return JSX(component, props, children...)
		}
	})
}

// CreateContext is used when you want to pass data to a deeply
// embedded child component without using props.
//
// See: https://reactjs.org/docs/context.html#reactcreatecontext
func CreateContext(defaultValue ...interface{}) (Context *js.Object, Provider *js.Object, Consumer *js.Object) {

	var res *js.Object

	if len(defaultValue) > 0 {
		res = React.Call("createContext", defaultValue[0])
	} else {
		res = React.Call("createContext")
	}

	return res, res.Get("Provider"), res.Get("Consumer")
}

// CloneElement is used to clone and return a new React Element.
//
// See: https://reactjs.org/docs/react-api.html#cloneelement
func CloneElement(element interface{}, props interface{}, children ...interface{}) *js.Object {

	args := []interface{}{
		element,
		SToMap(props),
	}
	if len(children) > 0 {
		args = append(args, children...)
	}

	return React.Call("cloneElement", args...)
}
