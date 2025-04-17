import * as pdfjsLib from 'pdfjs-dist';
pdfjsLib.GlobalWorkerOptions.workerSrc = `https://cdnjs.cloudflare.com/ajax/libs/pdf.js/5.0.375/pdf.worker.min.mjs`;

/**
 * File service for handling files
 */
const fileService = {
  /**
   * Extract text content from a file
   * @param {File} file - File object
   * @returns {Promise<string>} Extracted text content
   */
  extractContent: async (file) => {
    if (!file) return null;

    try {
      if (file.type === 'application/pdf' || file.name.endsWith('.pdf')) {
        return await extractPdfContent(file);
      } else if (file.type === 'text/plain' || file.name.endsWith('.txt')) {
        return await extractTextContent(file);
      } else if (file.type.includes('image/') || file.name.match(/\.(jpg|jpeg|png)$/i)) {
        return `[Image: ${file.name}]`;
      } else {
        console.warn(`Unsupported file type: ${file.type}`);
        return null;
      }
    } catch (error) {
      console.error(`Error extracting content from ${file.name}:`, error);
      return `[Error extracting content from ${file.name}]`;
    }
  },

  /**
   * Format file size
   * @param {number} size - File size in bytes
   * @returns {string} Formatted file size
   */
  formatFileSize: (size) => {
    if (size < 1024) {
      return `${size} B`;
    } else if (size < 1024 * 1024) {
      return `${(size / 1024).toFixed(1)} KB`;
    } else {
      return `${(size / (1024 * 1024)).toFixed(1)} MB`;
    }
  },
  
  /**
   * Convert an image file to base64
   * @param {File} file - Image file
   * @returns {Promise<string>} Base64 encoded image data
   */
  imageToBase64: (file) => {
    return new Promise((resolve, reject) => {
      if (!file || !file.type.includes('image/')) {
        reject(new Error('Invalid image file'));
        return;
      }
      
      const reader = new FileReader();
      
      reader.onload = (e) => {
        const base64String = e.target.result.split(',')[1];
        resolve(base64String);
      };
      
      reader.onerror = (error) => {
        reject(error);
      };
      
      reader.readAsDataURL(file);
    });
  }
};

/**
 * Extract text content from a text file
 * @param {File} file - Text file
 * @returns {Promise<string>} Text content
 */
const extractTextContent = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    
    reader.onload = (e) => {
      resolve(e.target.result);
    };
    
    reader.onerror = (error) => {
      reject(error);
    };
    
    reader.readAsText(file);
  });
};

/**
 * Extract text content from a PDF file
 * @param {File} file - PDF file
 * @returns {Promise<string>} Text content
 */
const extractPdfContent = async (file) => {
  try {
    const arrayBuffer = await file.arrayBuffer();
    const pdf = await pdfjsLib.getDocument({ data: arrayBuffer }).promise;
    let text = '';
    
    for (let i = 1; i <= pdf.numPages; i++) {
      const page = await pdf.getPage(i);
      const content = await page.getTextContent();
      const textItems = content.items.map(item => item.str);
      text += textItems.join(' ') + '\n\n';
    }
    
    return text.trim();
  } catch (error) {
    console.error('Error extracting PDF content:', error);
    throw error;
  }
};

export default fileService; 