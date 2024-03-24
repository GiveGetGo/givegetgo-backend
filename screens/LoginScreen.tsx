import React, { useState } from 'react';
import { StyleSheet, Text, View, TextInput, TouchableOpacity } from 'react-native';
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

  return (
    <View style={styles.container}>
      <Text style={styles.welcome}>Welcome!</Text>
      <TextInput
        style={styles.input}
        placeholder="Email"
        keyboardType="email-address"
        onChangeText={setEmail} // Update your state with the new text
      />
      <TextInput
        style={styles.input}
        placeholder="Password"
        secureTextEntry
        onChangeText={setPassword} // Update your state with the new text
      />
      <TouchableOpacity
        style={styles.loginButton}
        onPress={() => {
          navigation.navigate('MainScreen')
          console.log(email, password);
        }}
      >
        <Text style={styles.loginButtonText}>Login</Text>
      </TouchableOpacity>

      <TouchableOpacity onPress={() => navigation.navigate('ForgotPasswordScreen')}>             
        <Text style={styles.linkText}>Forget password?</Text>
      </TouchableOpacity>

      <TouchableOpacity onPress={() => navigation.navigate('SignUpScreen')}>             
        <Text style={styles.linkText}>Don’t have an account?</Text>
      </TouchableOpacity>

    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
    backgroundColor: '#fff', // adjust your background color
  },
  welcome: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 48,
  },
  input: {
    height: 40,
    width: '100%',
    marginVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#000', // adjust the color to match the design
    fontSize: 16,
    padding: 10,
  },
  loginButton: {
    backgroundColor: '#d3d3d3', // adjust the color to match the design
    paddingVertical: 12,
    paddingHorizontal: 20,
    borderRadius: 4,
    marginTop: 24,
  },
  loginButtonText: {
    color: '#000', // adjust the color to match the design
    fontSize: 16,
    fontWeight: '500',
  },
  forgotPassword: {
    marginTop: 12,
    color: '#000', // adjust the color to match the design
  },
  signUp: {
    marginTop: 12,
    color: '#000', // adjust the color to match the design
  },
  linkText: {
    color: 'blue',
    marginTop: 16,
  },
});

export default LoginScreen;