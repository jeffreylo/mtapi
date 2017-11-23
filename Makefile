.PHONY: proto

proto:
	protowrap -I proto proto/gtfs-realtime.proto proto/nyct-subway.proto
