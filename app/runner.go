package app

import (
  "context"
  "log/slog"
)

type Job struct {
  Ctx context.Context
  Origin string
  Destination string
  Data []byte

  // Tags []string
}

type Runner interface {

  Load(dir string, rootFile string) error
  SubmitJob(job *Job) error
}

type yaegiRunner struct {}

func (r *yaegiRunner) Load (dir string, rootFile string) error {
  slog.Debug("load runner directory", "dir", dir, "rootFile", rootFile)


  // add dir to artificial gopath for yaegi

  // import rootFile and have the imports of rootFile
  // handle user's requirements

  // Prepare yaegi to execute

  return nil
}

func (r *yaegiRunner) SubmitJob(job *Job) error {
  slog.Debug("submit job to runner")

  // TODO: We should add the context here that terminated the thread if it exceeds the 
  //        maximum time from the server.cfg file

  go r.processJob(job)

  return nil
}

func (r *yaegiRunner) processJob(job *Job) {

  slog.Info("processing job", "from", job.Origin, "to", job.Destination)
}


