import React from 'react';
import { View, StyleSheet, TouchableOpacity, SafeAreaView } from 'react-native';
import { Avatar, Button, Text, Card, Title, Paragraph, Appbar } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { useRoute } from '@react-navigation/native';

// Define the types for your navigation stack
type RootStackParamList = {
  HomeScreen: undefined;
  PostDetailsScreen: undefined;
  PostRequestInfoScreen: undefined;
  ProfileScreen: undefined;
};

type HomeScreenProps = NativeStackScreenProps<RootStackParamList, 'HomeScreen'>;

const PostDetailsScreen: React.FC<HomeScreenProps> = ({ navigation }: HomeScreenProps) => {
  const route = useRoute();
  const { postId } = route.params as { postId: string };

  console.log(postId) // use the postId to fetch corresponding info

  const goToRequestProfileScreen = () => {
    navigation.navigate('ProfileScreen');
  };

  const goToRequestInfoScreen = () => {
    navigation.navigate('PostRequestInfoScreen');
  };

  const use_navigation = useNavigation(); //for Appbar.BackAction

  return (
    <SafeAreaView  style={styles.container}> 
      <View style={styles.headerContainer}>
        <Appbar.BackAction style={styles.backAction} onPress={() => use_navigation.goBack()} />
        <Text style={styles.header}>GiveGetGo</Text>
        <View style={styles.backActionPlaceholder} />
      </View>
      <Card style={styles.card}>
        <Card.Content>
          <View style={styles.avatarContainer}>
            <TouchableOpacity onPress={goToRequestProfileScreen}>
                <Avatar.Image size={60} source={require('./profile_icon.jpg')} />
            </TouchableOpacity>
          </View>
          <Title style={styles.title}>Jimmy Ho</Title>
          <Title style={styles.boldText}>XXXXXX</Title>
          <Paragraph style={styles.paragraph}>XXXXXXXXXXXXXXXXXXXXXXXXXX</Paragraph>
        </Card.Content>
        <Card.Actions style={styles.cardActions}>
          <Button style={styles.button} mode="contained" onPress={goToRequestInfoScreen}>
            Send Request
          </Button>
        </Card.Actions>
      </Card>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 50,
    alignItems: 'center',
  },
  header: {
    fontSize: 20,
    fontWeight: 'bold',
    padding: 16,
    alignItems: 'center',
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    paddingLeft: 10, // Adds padding to the left of the avatar
    paddingRight: 10, // Adds padding to the right side
  },
  backActionPlaceholder: {
    width: 48, // This should match the width of the Appbar.BackAction for balance
    height: 48,
  },
  backAction: {
    marginLeft: 0 //This means the relative margin, comparing to the container (?)
  },
  card: {
    width: '100%',
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  avatarContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    marginVertical: 10,
  },
  title: {
    textAlign: 'center',
    marginVertical: 3,
  },
  boldText: {
    textAlign: 'center',
    fontWeight: 'bold',
    fontSize: 16, 
    marginVertical: 3,
  },
  paragraph: {
    textAlign: 'center',
    fontSize: 14,
    marginBottom: 20,
  },
  button: {
    // textAlign: 'center',
    // marginBottom: 10,
    position: 'absolute', 
    left: 40,
    right: 40, //position, left, right together controls the button's length and horizontal location
    alignSelf: 'center', 
  },
  cardActions: {
    justifyContent: 'center', 
    alignItems: 'center',
    padding: 15,
  },
});

export default PostDetailsScreen;
