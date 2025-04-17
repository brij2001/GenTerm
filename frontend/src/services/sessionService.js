import axios from 'axios';

/**
 * Session service for managing chat sessions
 */
const sessionService = {
  /**
   * Create a new chat session
   * @returns {Promise<string>} Session ID
   */
  createSession: async () => {
    try {
      const response = await axios.post('/api/session', {
        action: 'create'
      });
      
      return response.data.id;
    } catch (error) {
      console.error('Error creating session:', error);
      throw new Error('Failed to create session');
    }
  },

  /**
   * Get session details
   * @param {string} sessionId - Session ID
   * @returns {Promise<object>} Session details
   */
  getSession: async (sessionId) => {
    try {
      const response = await axios.post('/api/session', {
        action: 'get',
        id: sessionId
      });
      
      return response.data;
    } catch (error) {
      console.error('Error getting session:', error);
      throw new Error('Failed to get session');
    }
  }
};

export default sessionService; 