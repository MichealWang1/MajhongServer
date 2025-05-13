package service

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseServicePathKey(path string, serviceType uint16, serviceId uint16) string {
	return fmt.Sprintf("%s_%d_%d", path, serviceType, serviceId)
}

func ParseServicePath(serviceType uint16, serviceId uint16) string {
	return fmt.Sprintf("%d_%d", serviceType, serviceId)
}

func GetServiceType(servicePath string) uint16 {
	if len(servicePath) <= 0 {
		return 0
	}

	list := strings.Split(servicePath, "_")
	if len(list) < 2 {
		return 0
	}

	sType, err := strconv.Atoi(list[0])
	if err != nil {
		return 0
	}

	return uint16(sType)
}

func GetServiceId(servicePath string) uint16 {
	if len(servicePath) <= 0 {
		return 0
	}

	list := strings.Split(servicePath, "_")
	if len(list) < 2 {
		return 0
	}

	sId, err := strconv.Atoi(list[1])
	if err != nil {
		return 0
	}

	return uint16(sId)
}
