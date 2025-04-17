import React, { useRef, useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faUpload, faSpinner } from '@fortawesome/free-solid-svg-icons';

const FileUploader = ({ onFilesUploaded }) => {
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef(null);

  const handleDragOver = (e) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setIsDragging(false);
    
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      handleFiles(e.dataTransfer.files);
    }
  };

  const handleFileInputChange = (e) => {
    if (e.target.files && e.target.files.length > 0) {
      handleFiles(e.target.files);
    }
  };

  const handleFiles = async (fileList) => {
    setIsUploading(true);
    
    try {
      const files = Array.from(fileList);
      const validFiles = files.filter(file => {
        const fileType = file.type;
        return (
          fileType.includes('image/') || 
          fileType === 'application/pdf' || 
          fileType === 'text/plain' ||
          file.name.endsWith('.txt') ||
          file.name.endsWith('.pdf') ||
          file.name.endsWith('.png') ||
          file.name.endsWith('.jpg') ||
          file.name.endsWith('.jpeg')
        );
      });

      if (validFiles.length > 0) {
        onFilesUploaded(validFiles);
      }
    } catch (error) {
      console.error('Error uploading files:', error);
    } finally {
      setIsUploading(false);
    }
  };

  const handleButtonClick = () => {
    fileInputRef.current.click();
  };

  return (
    <div 
      className={`file-uploader ${isDragging ? 'dragging' : ''}`}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      <h3>Upload Files</h3>
      <p>Drag & drop files here or click to browse</p>
      <p>Supported formats: .txt, .pdf</p>
      <p>Unstable: .png, .jpg</p>
      <p>Note: working on fixing image support</p>
      <button 
        className="upload-button" 
        onClick={handleButtonClick}
        disabled={isUploading}
      >
        {isUploading ? (
          <>
            <FontAwesomeIcon icon={faSpinner} spin /> Processing...
          </>
        ) : (
          <>
            <FontAwesomeIcon icon={faUpload} /> Upload Files
          </>
        )}
      </button>
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileInputChange}
        style={{ display: 'none' }}
        multiple
        accept=".txt,.pdf,.png,.jpg,.jpeg"
      />
    </div>
  );
};

export default FileUploader; 