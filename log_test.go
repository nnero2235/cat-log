package catlog

import "testing"

func TestLog(t *testing.T){
	Debug("Debug-> %s","haha")
	Trace("trace-> %s","haha")
	Info("info-> %s","haha")
	Error("error-> %s","haha")
}

func TestConfig(t *testing.T){
	SetOptions(
		WithConsoleOutput(true),
		WithAsync(true),
		WithFileOutput(true),
		WithLevel(DEBUG),
	)
	Debug("Debug-> %s","haha")
	Trace("trace-> %s","haha")
	Info("info-> %s","haha")
	Error("error-> %s","haha")
}
