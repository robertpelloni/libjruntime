#pragma once

#include <cstdint>
#include <cstdio>
#include <string>
#include <sys/stat.h>
#include <unistd.h>
#include "Object.h"
#include "String.h"
#include "Array.h"
#include "Exceptions.h"
#include "GC.h"

namespace java {
namespace io {

class File : public ::java::lang::Object {
private:
    ::java::lang::String* path;

public:
    File(::java::lang::String* pathname) : path(pathname) {}

    virtual ::java::lang::String* getPath() const {
        return path;
    }

    virtual ::java::lang::String* getAbsolutePath() const {
        return path;
    }

    virtual bool exists() const {
        if (path == nullptr) return false;
        struct stat st;
        return (stat(path->getStdString().c_str(), &st) == 0);
    }

    virtual int64_t length() const {
        if (path == nullptr) return 0;
        struct stat st;
        if (stat(path->getStdString().c_str(), &st) == 0) {
            return static_cast<int64_t>(st.st_size);
        }
        return 0;
    }

    virtual bool deleteFile() {
        if (path == nullptr) return false;
        return (unlink(path->getStdString().c_str()) == 0);
    }
};

class FileInputStream : public ::java::lang::Object {
private:
    FILE* file;

public:
    FileInputStream(::java::lang::String* name) : file(nullptr) {
        if (name == nullptr) {
            throw ::java::lang::NullPointerException("File name is null");
        }
        file = fopen(name->getStdString().c_str(), "rb");
        if (file == nullptr) {
            throw ::java::lang::Exception("FileNotFoundException: " + name->getStdString());
        }
    }

    FileInputStream(File* f) : file(nullptr) {
        if (f == nullptr || f->getPath() == nullptr) {
            throw ::java::lang::NullPointerException("File is null");
        }
        file = fopen(f->getPath()->getStdString().c_str(), "rb");
        if (file == nullptr) {
            throw ::java::lang::Exception("FileNotFoundException: " + f->getPath()->getStdString());
        }
    }

    virtual ~FileInputStream() {
        close();
    }

    virtual int32_t read() {
        if (file == nullptr) {
            throw ::java::lang::Exception("IOException: Stream closed");
        }
        int c = fgetc(file);
        if (c == EOF) {
            return -1;
        }
        return static_cast<int32_t>(c & 0xFF);
    }

    virtual int32_t read(::java::lang::Array<int8_t>* buffer) {
        if (buffer == nullptr) {
            throw ::java::lang::NullPointerException("Buffer is null");
        }
        if (file == nullptr) {
            throw ::java::lang::Exception("IOException: Stream closed");
        }
        size_t n = fread(buffer->begin(), 1, buffer->length, file);
        if (n == 0 && feof(file)) {
            return -1;
        }
        return static_cast<int32_t>(n);
    }

    virtual int32_t read(::java::lang::Array<int32_t>* buffer) {
        if (buffer == nullptr) {
            throw ::java::lang::NullPointerException("Buffer is null");
        }
        if (file == nullptr) {
            throw ::java::lang::Exception("IOException: Stream closed");
        }
        int32_t count = 0;
        for (int32_t i = 0; i < buffer->length; ++i) {
            int c = fgetc(file);
            if (c == EOF) {
                if (count == 0) return -1;
                break;
            }
            (*buffer)[i] = static_cast<int32_t>(c & 0xFF);
            count++;
        }
        return count;
    }

    virtual void close() {
        if (file != nullptr) {
            fclose(file);
            file = nullptr;
        }
    }
};

class FileOutputStream : public ::java::lang::Object {
private:
    FILE* file;

public:
    FileOutputStream(::java::lang::String* name) : file(nullptr) {
        if (name == nullptr) {
            throw ::java::lang::NullPointerException("File name is null");
        }
        file = fopen(name->getStdString().c_str(), "wb");
        if (file == nullptr) {
            throw ::java::lang::Exception("IOException: Cannot open file for writing: " + name->getStdString());
        }
    }

    FileOutputStream(File* f) : file(nullptr) {
        if (f == nullptr || f->getPath() == nullptr) {
            throw ::java::lang::NullPointerException("File is null");
        }
        file = fopen(f->getPath()->getStdString().c_str(), "wb");
        if (file == nullptr) {
            throw ::java::lang::Exception("IOException: Cannot open file for writing: " + f->getPath()->getStdString());
        }
    }

    virtual ~FileOutputStream() {
        close();
    }

    virtual void write(int32_t b) {
        if (file == nullptr) {
            throw ::java::lang::Exception("IOException: Stream closed");
        }
        fputc(b & 0xFF, file);
    }

    virtual void write(::java::lang::Array<int8_t>* buffer) {
        if (buffer == nullptr) {
            throw ::java::lang::NullPointerException("Buffer is null");
        }
        if (file == nullptr) {
            throw ::java::lang::Exception("IOException: Stream closed");
        }
        fwrite(buffer->begin(), 1, buffer->length, file);
    }

    virtual void write(::java::lang::Array<int32_t>* buffer) {
        if (buffer == nullptr) {
            throw ::java::lang::NullPointerException("Buffer is null");
        }
        if (file == nullptr) {
            throw ::java::lang::Exception("IOException: Stream closed");
        }
        for (int32_t i = 0; i < buffer->length; ++i) {
            fputc((*buffer)[i] & 0xFF, file);
        }
    }

    virtual void flush() {
        if (file != nullptr) {
            fflush(file);
        }
    }

    virtual void close() {
        if (file != nullptr) {
            fclose(file);
            file = nullptr;
        }
    }
};

} // namespace io
} // namespace java
