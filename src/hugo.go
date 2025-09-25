package main

type HugoCmd struct {
	CommonConvert
}

func (hc *HugoCmd) Run(ctx *Context) error {
	defer hc.Input.Close()

	doz, _, entry, _, body, err := hc.getStuff(ctx)

	err = hc.GotoOutDir()
	if err != nil {
		return err
	}

	err = hc.outBody(body)

	err = hc.outPhotos(doz, entry)
	if err != nil {
		return err
	}

	return nil
}
