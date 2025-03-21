package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock dependencies
type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) DoesExist(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func (m *MockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	args := m.Called(filename, data, perm)
	return args.Error(0)
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *MockFileSystem) BackupFile(src, dst string) error {
	args := m.Called(src, dst)
	return args.Error(0)
}

// Mock logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Fatal(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.Called(format, args)
}

// Test function
func TestCreateConfig(t *testing.T) {
	// Setup
	origDoesExist := doesExist
	origWriteFile := os.WriteFile
	origMkdirAll := os.MkdirAll
	origBackupFile := backupFile
	origLogger := Logger

	// Create mocks
	mockFS := new(MockFileSystem)
	mockLog := new(MockLogger)

	// Replace global functions with mocks
	doesExist = mockFS.DoesExist
	os.WriteFile = mockFS.WriteFile
	os.MkdirAll = mockFS.MkdirAll
	backupFile = mockFS.BackupFile
	Logger = mockLog

	// Cleanup after test
	defer func() {
		doesExist = origDoesExist
		os.WriteFile = origWriteFile
		os.MkdirAll = origMkdirAll
		backupFile = origBackupFile
		Logger = origLogger
	}()

	// Test data
	cfgBinary := "test-binary"
	cfgFile := "/path/to/config/.tp.toml"
	cfgMdFile := "test.md"
	cfgPlanFile := "test.plan"
	configDir := filepath.Dir(cfgFile)
	mockConfig := []byte("test config content")

	// Test Case 1: Create new config file in existing directory
	t.Run("CreateNewConfigInExistingDir", func(t *testing.T) {
		// Reset expectations
		mockFS = new(MockFileSystem)
		mockLog = new(MockLogger)

		// Setup expectations
		// Mock createOrOverwrite to return configExists=false, createFile=true
		createOrOverwrite = func(string) (bool, bool) {
			return false, true
		}

		// Mock genConfig to return test config
		genConfig = func(ConfigParams) ([]byte, error) {
			return mockConfig, nil
		}

		mockFS.On("DoesExist", configDir).Return(true)
		mockFS.On("WriteFile", cfgFile, mockConfig, os.FileMode(0o600)).Return(nil)
		mockLog.On("Debugf", "Inside configExists and 'config' is: %s", []interface{}{string(mockConfig)})
		mockLog.On("Infof", "Config file %s created", []interface{}{cfgFile})

		// Call function under test
		err := createConfig(cfgBinary, cfgFile, cfgMdFile, cfgPlanFile)

		// Assertions
		assert.NoError(t, err)
		mockFS.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	})

	// Test Case 2: Create new config file in non-existing directory
	t.Run("CreateNewConfigInNewDir", func(t *testing.T) {
		// Reset expectations
		mockFS = new(MockFileSystem)
		mockLog = new(MockLogger)

		// Setup expectations
		// Mock createOrOverwrite to return configExists=false, createFile=true
		createOrOverwrite = func(string) (bool, bool) {
			return false, true
		}

		// Mock genConfig to return test config
		genConfig = func(ConfigParams) ([]byte, error) {
			return mockConfig, nil
		}

		mockFS.On("DoesExist", configDir).Return(false)
		mockFS.On("MkdirAll", configDir, os.FileMode(0o750)).Return(nil)
		mockFS.On("WriteFile", cfgFile, mockConfig, os.FileMode(0o600)).Return(nil)
		mockLog.On("Debugf", "Inside configExists and 'config' is: %s", []interface{}{string(mockConfig)})
		mockLog.On("Infof", "Config file %s created", []interface{}{cfgFile})

		// Call function under test
		err := createConfig(cfgBinary, cfgFile, cfgMdFile, cfgPlanFile)

		// Assertions
		assert.NoError(t, err)
		mockFS.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	})

	// Test Case 3: Overwrite existing config file
	t.Run("OverwriteExistingConfig", func(t *testing.T) {
		// Reset expectations
		mockFS = new(MockFileSystem)
		mockLog = new(MockLogger)

		// Setup expectations
		// Mock createOrOverwrite to return configExists=true, createFile=true
		createOrOverwrite = func(string) (bool, bool) {
			return true, true
		}

		// Mock genConfig to return test config
		genConfig = func(ConfigParams) ([]byte, error) {
			return mockConfig, nil
		}

		// Mock time.Now()
		timeFormat := time.Now().Local().Format("200601021504")
		backupFile := cfgFile + "-" + timeFormat

		mockFS.On("DoesExist", configDir).Return(true)
		mockFS.On("BackupFile", cfgFile, backupFile).Return(nil)
		mockFS.On("WriteFile", cfgFile, mockConfig, os.FileMode(0o600)).Return(nil)
		mockLog.On("Debugf", "Config is: \n%s\n", []interface{}{string(mockConfig)})
		mockLog.On("Infof", "Backup file %s created", []interface{}{backupFile})
		mockLog.On("Infof", "Config file %s created", []interface{}{cfgFile})

		// Call function under test
		err := createConfig(cfgBinary, cfgFile, cfgMdFile, cfgPlanFile)

		// Assertions
		assert.NoError(t, err)
		mockFS.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	})

	// Test Case 4: Print config only (don't create file)
	t.Run("PrintConfigOnly", func(t *testing.T) {
		// Reset expectations
		mockFS = new(MockFileSystem)
		mockLog = new(MockLogger)

		// Setup expectations
		// Mock createOrOverwrite to return configExists=false, createFile=false
		createOrOverwrite = func(string) (bool, bool) {
			return false, false
		}

		// Mock genConfig to return test config
		genConfig = func(ConfigParams) ([]byte, error) {
			return mockConfig, nil
		}

		mockLog.On("Info", string(mockConfig))

		// Call function under test
		err := createConfig(cfgBinary, cfgFile, cfgMdFile, cfgPlanFile)

		// Assertions
		assert.NoError(t, err)
		mockFS.AssertExpectations(t)
		mockLog.AssertExpectations(t)
	})
}
