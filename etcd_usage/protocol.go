package main

type Job struct {
	Name     string   `json:"name"`
	Command  string   `json:"command"`
	CronExpr string   `json:"cronExpr"`
}


