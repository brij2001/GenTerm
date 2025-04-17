import axios from 'axios';

/**
 * Chat service for interacting with LLM API
 */
const chatService = {
  /**
   * Send a query to the LLM with context
   * @param {string} sessionId - Session ID
   * @param {string} query - User query
   * @param {string} context - Content from files
   * @returns {Promise<string>} LLM response
   */
  sendQuery: async (sessionId, query, context) => {
    try {
      const contextArray = context ? [context] : [];
      
      const response = await axios.post('/api/chat', {
        sessionId,
        query,
        context: contextArray
      });
      
      return response.data.response;
    } catch (error) {
      console.error('Error sending query:', error);
      throw new Error('Failed to get AI response');
    }
  },
  
  /**
   * Send a query with an image to the LLM
   * @param {string} sessionId - Session ID
   * @param {string} query - User query
   * @param {string} base64Image - Base64 encoded image data
   * @param {string} context - Content from files
   * @returns {Promise<string>} LLM response
   */
  sendImageQuery: async (sessionId, query, base64Image, context) => {
    try {
      const contextArray = context ? [context] : [];
      
      const messageContent = [
        { type: "text", text: query },
        {
          type: "image_url",
          image_url: {
            url: `data:image/jpeg;base64,${base64Image}`
          }
        },
      ];
      
      const response = await axios.post('/api/chat', {
        sessionId,
        query,
        context: contextArray,
        messageContent
      });
      
      return response.data.response;
    } catch (error) {
      console.error('Error sending image query:', error);
      throw new Error('Failed to get AI response for image');
    }
  }
};

export default chatService; 