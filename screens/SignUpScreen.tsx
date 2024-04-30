import React, { useState } from 'react';
import { StyleSheet, View, Text, TextInput, TouchableOpacity } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

type RootStackParamList = {
  SignUpScreen: undefined;
  CheckEmailScreen: undefined;
  LoginScreen: undefined;
};

type ScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'SignUpScreen' | 'CheckEmailScreen' | 'LoginScreen'
>;

const SignUpScreen: React.FC = () => {
  const navigation = useNavigation<ScreenNavigationProp>();
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [schoolClass, setschoolClass] = useState<string>(''); // 'class' is a reserved word
  const [major, setMajor] = useState<string>('');

  const handleSignUp = () => {
    registerUser(email, password, schoolClass, major);
    console.log('Signing up with:', email, password, schoolClass, major);          
  };

  const handleLogIn = () => {
    navigation.navigate('LoginScreen');
  };

  async function registerUser(email: string, password: string, schoolClass: string, major: string) {
    navigation.navigate('CheckEmailScreen'); //for testing
    try {
      const response = await fetch('http://api.givegetgo.xyz/v1/user/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: email,
          password: password,
          class: schoolClass,
          major: major
        }),
      });

      const json = await response.json(); // Parse the JSON response
      console.log("Registration info", json);

      if (response.status === 201) {
        console.log('Registration successful:', json);
        requestMFASetup();
        navigation.navigate('CheckEmailScreen');
      } else {
        // Handle different types of errors based on response status
        console.error('Registration failed:', json.msg);
        alert(`Registration failed: ${json.msg}`);
      }
    } catch (error) {
      // Handle network errors or other unexpected issues
      console.error('Network error:', error);
      alert('Failed to connect to the server. Please try again later.');
    }
  }

  async function requestMFASetup() {
    try {
      const response = await fetch('http://api.givegetgo.xyz/v1/mfa', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: email,
        }),
      });
  
      const json = await response.json(); // Parse the JSON response
      console.log("MFA setup info", json);
  
      if (response.status === 200) {
        console.log('MFA setup successful:', json);
        // Here you can navigate or perform further actions upon successful setup
      } else {
        // Handle different types of errors based on response status
        console.error('MFA setup failed:', json.msg);
        alert(`MFA setup failed: ${json.msg}`);
      }
    } catch (error) {
      // Handle network errors or other unexpected issues
      console.error('Network error:', error);
      alert('Failed to connect to the server for MFA setup. Please try again later.');
    }
  }

  return (
    <View style={styles.container}>
      <Text style={styles.header}>Hello!</Text>
      <Text style={styles.subHeader}>Let's create an account</Text>
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
        onChangeText={setPassword}
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
      />
      <TouchableOpacity style={styles.button} onPress={handleSignUp}>
        <Text style={styles.buttonText}>Sign Up</Text>
      </TouchableOpacity>
      <View style={styles.footer}>
        <Text style={styles.footerText}>
          Have an account?{' '}
          <Text style={styles.link} onPress={handleLogIn}>Log In</Text>
        </Text>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#fff',
  },
  header: {
    fontSize: 42,
    fontWeight: 'bold',
    marginTop: 160,
    marginBottom: 10,
    marginLeft: '10%',
    alignSelf: 'flex-start',
  },
  subHeader: {
    fontSize: 18,
    color: 'grey',
    marginBottom: 20,
    marginLeft: '10%',
    alignSelf: 'flex-start',
  },
  input: {
    width: '80%',
    borderBottomWidth: 1,
    borderBottomColor: 'grey',
    fontSize: 16,
    paddingVertical: 10,
    marginBottom: 20,
  },
  button: {
    width: '75%',
    backgroundColor: 'black',
    padding: 13,
    justifyContent: 'center',
    alignItems: 'center',
    borderRadius: 5,
    marginTop: 20,
  },
  buttonText: {
    color: '#FAFAFA',
    fontSize: 18,
    fontWeight: '500',
  },
  footer: {
    position: 'absolute',
    bottom: 50,
    alignSelf: 'center', // This will center the footer horizontally
  },
  footerText: {
    fontSize: 16,
    color: '#000',
    marginRight: 5,
  },
  link: {
    fontSize: 16,
    fontWeight: '600',
  },
});

export default SignUpScreen;


// import React, { useState } from 'react';
// import { StyleSheet, View, Text, TextInput, TouchableOpacity } from 'react-native';
// import { useNavigation } from '@react-navigation/native';
// import { StackNavigationProp } from '@react-navigation/stack';

// // Define the types for your navigation stack
// type RootStackParamList = {
//   SignUpScreen: undefined;
//   CheckEmailScreen: undefined;
//   LoginScreen: undefined;
// };

// // Define the type for the navigation prop
// type ScreenNavigationProp = StackNavigationProp<
//   RootStackParamList,
//   'SignUpScreen' | 'CheckEmailScreen' | 'LoginScreen'
// >;

// const SignUpScreen: React.FC = () => {
//   const [email, setEmail] = useState<string>('');
//   const [password, setPassword] = useState<string>('');
//   const [schoolClass, setschoolClass] = useState<string>(''); // 'class' is a reserved word
//   const [major, setMajor] = useState<string>('');
//   const navigation = useNavigation<ScreenNavigationProp>();

//   const handleSignUp = () => { 
//     // Handle the sign up logic
//     navigation.navigate('CheckEmailScreen');
//     console.log('Signing up with:', email, password, schoolClass, major);          
//     // Add validation for password match and call the API to sign up
//   };

//   const handleLogIn = () => { 
//     navigation.navigate('LoginScreen');
//   };

//   return (
//     <View style={styles.container}>
//       <Text style={styles.titleText}>Sign Up</Text>
      
//       <TextInput
//         style={styles.input}
//         placeholder="Email"
//         onChangeText={setEmail}
//         value={email}
//         keyboardType="email-address"
//         autoCapitalize="none"
//       />
      
//       <TextInput
//         style={styles.input}
//         placeholder="Password"
//         secureTextEntry
//         onChangeText={setPassword}
//         value={password}
//       />
      
//       <TextInput
//         style={styles.input}
//         placeholder="Class"
//         onChangeText={setschoolClass}
//         value={schoolClass}
//       />

//       <TextInput
//         style={styles.input}
//         placeholder="Major"
//         onChangeText={setMajor}
//         value={major}
//       />

//       <TouchableOpacity style={styles.button} onPress={handleSignUp}>
//         <Text style={styles.buttonText}>Sign Up</Text>
//       </TouchableOpacity>

//       <Text style={styles.text}>
//         Already have an account?{' '}
//         <TouchableOpacity onPress={handleLogIn}>
//           <Text style={styles.linkText}>Log In</Text>
//         </TouchableOpacity>
//       </Text>

//     </View>
//   );
// };

// const styles = StyleSheet.create({
//   container: {
//     flex: 1,
//     justifyContent: 'center',
//     alignItems: 'center',
//     padding: 16,
//     backgroundColor: '#fff',
//   },
//   titleText: {
//     fontSize: 24,
//     fontWeight: 'bold',
//     marginBottom: 48,
//   },
//   input: {
//     height: 40,
//     width: '100%',
//     borderColor: 'gray',
//     borderWidth: 1,
//     marginBottom: 16,
//     paddingHorizontal: 8,
//   },
//   button: {
//     backgroundColor: 'gray',
//     padding: 10,
//     borderRadius: 5,
//     width: '100%',
//     alignItems: 'center',
//   },
//   buttonText: {
//     color: '#fff',
//   },
//   linkText: {
//     color: 'blue',
//     marginTop: 16,
//   },
// });

// export default SignUpScreen;
