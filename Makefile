# Copyright © 2021 Tim Ebert
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

REPOSITORY := docker.pkg.github.com/timebertt/speedtest-exporter/speedtest-exporter
TAG := dev

test:
	go test ./...

docker-image:
	DOCKER_BUILDKIT=1 docker build -t $(REPOSITORY):$(TAG) --target speedtest-exporter .

docker-up:
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose up -d --build

docker-down:
	docker-compose down
