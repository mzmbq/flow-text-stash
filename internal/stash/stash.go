package stash

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/atotto/clipboard"
	"github.com/mzmbq/flow-launcher-go"
	"github.com/mzmbq/flow-text-stash/internal/store"
	"github.com/mzmbq/flow-text-stash/internal/utils"
)

type TextStash struct {
	flow.Plugin
	Store *store.Store
}

func New(s *store.Store) *TextStash {
	ts := &TextStash{
		Plugin: *flow.NewPlugin(),
		Store:  s,
	}

	ts.Query(ts.handleQuery)
	ts.Action("paste", ts.handlePaste)
	ts.Action("create", ts.handleCreate)
	ts.Action("open_data", ts.handleOpenData)
	ts.ContextMenu(ts.handleCtxMenu)

	return ts
}

func (ts *TextStash) listAllStashes(req *flow.Request) *flow.Response {
	res := flow.NewResponse(req)
	for k, v := range ts.Store.Data {
		res.AddResult(&flow.Result{
			Title:    k,
			SubTitle: utils.Wrap(v),
			IcoPath:  "pencil.png",
			RpcAction: &flow.JsonRpcAction{
				Method:     "paste",
				Parameters: []string{v},
			},
		})
	}
	return res
}

func (ts *TextStash) handleQuery(req *flow.Request) *flow.Response {
	res := flow.NewResponse(req)
	target := req.Parameters[0]
	if target == "" {
		return ts.listAllStashes(req)
	}
	matches := ts.Store.GetFuzzy(target)

	// List matches
	for _, m := range matches {
		val := ts.Store.Data[m]
		res.AddResult(&flow.Result{
			Title:    m,
			SubTitle: utils.Wrap(val),
			IcoPath:  "pencil.png",
			RpcAction: &flow.JsonRpcAction{
				Method:     "paste",
				Parameters: []string{val},
			},
		})
	}

	// List create option for when no matches or no exact match was found
	if len(matches) == 0 || matches[0] != target {
		res.AddResult(&flow.Result{
			Title:    fmt.Sprintf("Create a paste: %s", target),
			SubTitle: "",
			IcoPath:  "add.png",
			RpcAction: &flow.JsonRpcAction{
				Method:     "create",
				Parameters: []string{target},
			},
		})
	}

	return res
}

func (ts *TextStash) handleCreate(params []string) *flow.Response {
	text, err := clipboard.ReadAll()
	if err != nil {
		// TODO: handle
		return nil
	}

	ts.Store.Set(params[0], text)
	err = ts.Store.Save()
	if err != nil {
		// TODO: handle
		return nil
	}
	return nil
}

func (ts *TextStash) handlePaste(params []string) *flow.Response {
	err := clipboard.WriteAll(params[0])
	if err != nil {
		// TODO: handle
		return nil
	}
	return nil
}

func (ts *TextStash) handleOpenData([]string) *flow.Response {
	cmd := exec.Command("cmd", "/c", "start", ts.Store.Path)
	if err := cmd.Run(); err != nil {
		log.Printf("Error opening file: %v", err)
	}
	return nil
}

func (ts *TextStash) handleCtxMenu(req *flow.Request) *flow.Response {
	res := flow.NewResponse(req)
	res.AddResult(&flow.Result{
		Title:    "Edit stashes",
		SubTitle: "Open data.yaml",
		IcoPath:  "pencil.png",
		RpcAction: &flow.JsonRpcAction{
			Method:     "open_data",
			Parameters: make([]string, 0),
		},
	})

	return res
}
