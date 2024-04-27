import React, { useState } from 'react';
import { StyleSheet, Text, View, TextInput, TouchableOpacity } from 'react-native';
import { Button, Modal} from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

// Define the types for your navigation stack
type RootStackParamList = {
  ForgotPasswordScreen: undefined;
  SignUpScreen: undefined;
  MainScreen: undefined;
};

// Define the type for the navigation prop
type LoginScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'ForgotPasswordScreen' | 'SignUpScreen' | 'MainScreen'
>;

const LoginScreen: React.FC = () => {
  const navigation = useNavigation<LoginScreenNavigationProp>();
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [visible, setVisible] = useState(false); // For modal visibility
  const [error, setError] = useState(''); // To hold the error message

  const hideModal = () => setVisible(false);
  const showModal = (message: string) => {
    setError(message);
    setVisible(true);
  };
  
  const handleLogin = async () => {
    navigation.navigate('MainScreen'); 
    // showModal('Login failed: Incorrect email or password.'); //uncomment this to format modal
    // try {
    //   const response = await fetch('http://api.givegetgo.xyz/v1/user/login', {
    //     method: 'POST',
    //     headers: {
    //       'Content-Type': 'application/json',
    //     },
    //     body: JSON.stringify({
    //       email: email,
    //       password: password,
    //     }),
    //   });
  
    //   const json = await response.json();
    //   console.log("login info", json)
  
    //   if (response.status === 200) {
    //     console.log('Login successful:', json);
    //     navigation.navigate('MainScreen'); // Navigate to the main screen on success
    //   } else {
    //     // Handle different types of errors based on response status
    //     console.error('Login failed:', json.msg);
    //     alert(`Login failed: ${json.msg}`);
    //   }
    // } catch (error) {
    //   // Handle network errors or other unexpected issues
    //   console.error('Network error:', error);
    //   alert('Failed to connect to the server. Please try again later.');
    // }
  };

  return (
    <View style={styles.container}>
      <View style={styles.contentWrapper}>
        <Text style={styles.welcome}>Welcome!</Text>
        <TextInput
          style={styles.input}
          placeholder="Email"
          keyboardType="email-address"
          autoCapitalize="none"
          onChangeText={setEmail}
        />
        <TextInput
          style={styles.input}
          placeholder="Password"
          secureTextEntry
          autoCapitalize="none"
          onChangeText={setPassword}
        />
      </View>
      <TouchableOpacity
        style={styles.loginButton}
        onPress={() => handleLogin()}
      >
        <Text style={styles.loginButtonText}>login</Text>
      </TouchableOpacity>
        <Text 
          style={styles.linkText}
          onPress={() => navigation.navigate('ForgotPasswordScreen')}
        >
          Forget password?
        </Text>
        <Text 
          style={styles.linkText}
          onPress={() => navigation.navigate('SignUpScreen')}
        >
          Don’t have an account?
        </Text>
        <Modal visible={visible} onDismiss={hideModal} contentContainerStyle={styles.modalContainer}>
          <Text style={styles.errorText}>{error}</Text>
          <Button mode="contained" onPress={hideModal}>
            Close
          </Button>
        </Modal>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
    backgroundColor: '#fff',
  },
  contentWrapper: {
    alignItems: 'flex-start',
    padding: 20,
    width: '100%',
  },
  welcome: {
    fontSize: 42,
    fontWeight: 'bold',
    alignSelf: 'flex-start',
    marginTop: 0,
    marginBottom: 16,
  },
  input: {
    height: 48, // Increase the height for a better tap area
    fontSize: 16,
    width: '100%',
    marginVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#000',
    padding: 5,
  },
  loginButton: {
    backgroundColor: 'black',
    width: '80%', // Adjust the width as needed
    alignItems: 'center',
    paddingVertical: 12,
    borderRadius: 4,
    marginTop: 15,
    marginBottom: 21,
  },
  loginButtonText: {
    color: "#FAFAFA",
    fontSize: 16,
    fontWeight: '500',
  },
  linkText: {
    // color: 'blue',
    marginTop: 10,
    fontSize: 16,
  },
  modalContainer: {
    backgroundColor: '#FAFAFA',
    padding: 20,
    margin: 20,
    borderRadius: 5,
    alignItems: 'center',
  },
  errorText: {
    fontSize: 16,
    marginBottom: 20,
    textAlign: 'center',
  },
});

export default LoginScreen;
