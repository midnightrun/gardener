// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webhook

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

// Validator validates objects.
type Validator interface {
	Validate(ctx context.Context, new, old runtime.Object) error
}

type validationWrapper struct {
	Validator
}

// Mutate implements the `Mutator` interface and calls the `Validate` function of the underlying validator.
func (d *validationWrapper) Mutate(ctx context.Context, new, old runtime.Object) error {
	return d.Validate(ctx, new, old)
}

// InjectFunc calls the inject.Func on the handler mutators.
func (d *validationWrapper) InjectFunc(f inject.Func) error {
	if err := f(d.Validator); err != nil {
		return errors.Wrap(err, "could not inject into the validator")
	}
	return nil
}

func hybridValidator(val Validator) Mutator {
	return &validationWrapper{val}
}