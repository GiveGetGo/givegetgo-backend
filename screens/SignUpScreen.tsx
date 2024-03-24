import React, { useState } from 'react';
import { StyleSheet, View, Text, TextInput, TouchableOpacity } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

// Define the types for your navigation stack
type RootStackParamList = {
  SignUpScreen: undefined;
  CheckEmailScreen: undefined;
  LoginScreen: undefined;
};

// Define the type for the navigation prop
type ScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'SignUpScreen' | 'CheckEmailScreen' | 'LoginScreen'
>;

const SignUpScreen: React.FC = () => {
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [schoolClass, setschoolClass] = useState<string>(''); // 'class' is a reserved word
  const [major, setMajor] = useState<string>('');
  const navigation = useNavigation<ScreenNavigationProp>();

  const handleSignUp = () => { 
    // Handle the sign up logic
    navigation.navigate('CheckEmailScreen');
    console.log('Signing up with:', email, password, schoolClass, major);          
    // Add validation for password match and call the API to sign up
  };

  const handleLogIn = () => { 
    navigation.navigate('LoginScreen');
  };

  return (
    <View style={styles.container}>
      <Text style={styles.titleText}>Sign Up</Text>
      
      <TextInput
        style={styles.input}
        placeholder="Email"
        onChangeText={setEmail}
        value={email}
        keyboardType="email-address"
        autoCapitalize="none"
      />
      
      <TextInput
        style={styles.input}
        placeholder="Password"
        secureTextEntry
        onChangeText={setPassword}
        value={password}
      />
      
      <TextInput
        style={styles.input}
        placeholder="Class"
        onChangeText={setschoolClass}
        value={schoolClass}
      />

      <TextInput
        style={styles.input}
        placeholder="Major"
        onChangeText={setMajor}
        value={major}
      />

      <TouchableOpacity style={styles.button} onPress={handleSignUp}>
        <Text style={styles.buttonText}>Sign Up</Text>
      </TouchableOpacity>

      <Text style={styles.text}>
        Already have an account?{' '}
        <TouchableOpacity onPress={handleLogIn}>
          <Text style={styles.linkText}>Log In</Text>
        </TouchableOpacity>
      </Text>

    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 16,
    backgroundColor: '#fff',
  },
  titleText: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 48,
  },
  input: {
    height: 40,
    width: '100%',
    borderColor: 'gray',
    borderWidth: 1,
    marginBottom: 16,
    paddingHorizontal: 8,
  },
  button: {
    backgroundColor: 'gray',
    padding: 10,
    borderRadius: 5,
    width: '100%',
    alignItems: 'center',
  },
  buttonText: {
    color: '#fff',
  },
  linkText: {
    color: 'blue',
    marginTop: 16,
  },
});

export default SignUpScreen;
