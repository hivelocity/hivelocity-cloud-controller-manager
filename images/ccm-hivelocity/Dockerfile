# Copyright 2023 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build the manager binary
FROM --platform=${BUILDPLATFORM} docker.io/library/golang:1.21.4-bullseye@sha256:ea48852045248675a6f1948acc4c996b1a5de151a77213b936840bf1f82ddd02 as build
ARG TARGETOS TARGETARCH

COPY . /src/cluster-api-provider-hivelocity
WORKDIR /src/cluster-api-provider-hivelocity
RUN --mount=type=cache,target=/root/.cache --mount=type=cache,target=/go/pkg \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \
    go build -mod=readonly -ldflags "${LDFLAGS} -extldflags '-static'" \
    -o manager main.go

FROM --platform=${BUILDPLATFORM} gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /src/cluster-api-provider-hivelocity/manager .
# Use uid of nonroot user (65532) because kubernetes expects numeric user when applying pod security policies
USER 65532
ENTRYPOINT ["/manager"]
