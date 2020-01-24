/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"os"

	"github.com/buildpacks/libbuildpack/v2/buildplan"
	"github.com/cloudfoundry/google-stackdriver-cnb/java"
	"github.com/cloudfoundry/libcfbuildpack/v2/detect"
)

func main() {
	detect, err := detect.DefaultDetect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Detect: %s\n", err)
		os.Exit(101)
	}

	if code, err := d(detect); err != nil {
		detect.Logger.TerminalError(detect.Buildpack, err.Error())
		os.Exit(code)
	} else {
		os.Exit(code)
	}
}

func d(detect detect.Detect) (int, error) {
	d := detect.Services.HasService("google-stackdriver-debugger", "PrivateKeyData")
	p := detect.Services.HasService("google-stackdriver-profiler", "PrivateKeyData")

	if !d && !p {
		return detect.Fail(), nil
	}

	q := buildplan.Plan{
		Requires: []buildplan.Required{
			{Name: "jvm-application"},
		},
	}

	if d {
		q.Provides = append(q.Provides, buildplan.Provided{Name: java.DebuggerDependency})
		q.Requires = append(q.Requires, buildplan.Required{Name: java.DebuggerDependency})
	}

	if p {
		q.Provides = append(q.Provides, buildplan.Provided{Name: java.ProfilerDependency})
		q.Requires = append(q.Requires, buildplan.Required{Name: java.ProfilerDependency})
	}

	return detect.Pass(q)
}
