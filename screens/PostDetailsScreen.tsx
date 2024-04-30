import React, { useEffect, useState } from 'react';
import { View, StyleSheet, TouchableOpacity, SafeAreaView } from 'react-native';
import { Avatar, Button, Text, Card, Title, Paragraph, Appbar } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';
import { useRoute } from '@react-navigation/native';
import { useFonts, Montserrat_700Bold_Italic } from '@expo-google-fonts/montserrat';

// Define the types for your navigation stack
type RootStackParamList = {
  HomeScreen: undefined;
  PostDetailsScreen: undefined;
  PostRequestInfoScreen: { postId: string, name: string };
  ProfileScreen: undefined;
};

type HomeScreenProps = NativeStackScreenProps<RootStackParamList, 'HomeScreen'>;

type PostOwnerProfile = { 
  post_id: string;
  name: string;
  title: string;
  description: string;
  avatar: string;
}

const defaultPostOwnerProfile: PostOwnerProfile = {
  post_id: '1',
  name: 'Kevin Lu',
  title: 'Taking Care of My Cat',
  description: 'I will be out of town on February 28th and need someone to look over my cat Lana.',
  avatar: '../assets/avatars/avatar5.png',
};

const PostDetailsScreen: React.FC<HomeScreenProps> = ({ navigation }: HomeScreenProps) => {

  const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic });

  const route = useRoute();
  const { postId } = route.params as { postId: string };

  console.log(postId) // use the postId to fetch corresponding info

  const [postOwnerProfile, setPostOwnerProfile] = useState<PostOwnerProfile>(defaultPostOwnerProfile);

  useEffect(() => {                                                                                     //fill this in to get db info //use PostID to fetch the corresponding profile 
    const fetchPostOwnerProfile = async () => {
      try {
        const response = await fetch('URL_TO_YOUR_BACKEND/bid/{postId}');
        const json = await response.json();
        setPostOwnerProfile(json); // Adjust this depending on the structure of your JSON
      } catch (error) {
        // console.error(error); // uncomment this after finish frontend developing
      }
    };
    fetchPostOwnerProfile();
  }, []);

  const goToRequestProfileScreen = () => {
    navigation.navigate('ProfileScreen');
  };

  const goToRequestInfoScreen = (postId: string, name: string) => {
    navigation.navigate('PostRequestInfoScreen', {postId: postId, name: name});
  };

  const use_navigation = useNavigation(); //for Appbar.BackAction


  // require() could not directly take the url
  interface ImageMap {
    [key: string]: NodeRequire;
  }
  const imageMap: { [key: string]: NodeRequire } = {
      '../assets/avatars/avatar1.png': require('../assets/avatars/avatar1.png'),
      '../assets/avatars/avatar2.png': require('../assets/avatars/avatar2.png'),
      '../assets/avatars/avatar3.png': require('../assets/avatars/avatar3.png'),
      '../assets/avatars/avatar4.png': require('../assets/avatars/avatar4.png'),
      '../assets/avatars/avatar5.png': require('../assets/avatars/avatar5.png'),
      '../assets/avatars/avatar6.png': require('../assets/avatars/avatar6.png'),
      '../assets/avatars/avatar7.png': require('../assets/avatars/avatar7.png'),
      '../assets/avatars/avatar8.png': require('../assets/avatars/avatar8.png'),
      '../assets/avatars/avatar9.png': require('../assets/avatars/avatar9.png'),
  };
  const getProfilePictureSource = (uri: string) => {
      return imageMap[uri] || require(`./profile_icon.jpg`);
  };


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
                <Avatar.Image size={70} source={getProfilePictureSource(postOwnerProfile.avatar)} />
            </TouchableOpacity>
          </View>
          <Title style={styles.title}>{postOwnerProfile.name}</Title>
          <Title style={styles.boldText}>{postOwnerProfile.title}</Title>
          <Paragraph style={styles.paragraph}>{postOwnerProfile.description}</Paragraph>
        </Card.Content>
        <Card.Actions style={styles.cardActions}>
          <Button
            style={styles.button}
            mode="contained"
            onPress={() => goToRequestInfoScreen(postId, postOwnerProfile.name)} 
          >
            Send Bid
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
    justifyContent: 'center',
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    justifyContent: 'space-between', // Distributes items evenly horizontally
    paddingLeft: 10, 
    paddingRight: 10, 
    position: 'absolute', // So that while setting card to the vertical middle, it still stays at the same place
    top: 0, 
    left: 0,
    right: 0,
    zIndex: 1, // Ensure the headerContainer is above the card
  },
  header: {
    fontSize: 22, // Increase the font size
    fontWeight: '600', // Make the font weight bold
    fontFamily: 'Montserrat_700Bold_Italic',
    textAlign: 'center', // Center the text
    color: '#444444', // Dark gray color
  },
  backActionPlaceholder: {
    width: 48, // This should match the width of the Appbar.BackAction for balance
    height: 48,
  },
  backAction: {
    marginLeft: 0 //This means the relative margin, comparing to the container (?)
  },
  card: { //page gets longer when there are more contexts
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 6,
    marginHorizontal: 12,
    elevation: 0, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 15, // Add padding inside the card
    // marginTop: 170,
  },
  avatarContainer: {
    alignItems: 'center',
    justifyContent: 'center',                    
    marginBottom: 0,
    marginTop: -7,
  },
  title: {
    textAlign: 'center',
    marginBottom: -5,
  },
  boldText: {
    textAlign: 'center',
    fontWeight: 'bold',
    fontSize: 16, 
    marginBottom: -2,
  },
  paragraph: {
    textAlign: 'center',
    fontSize: 14,
    marginVertical: 0,
    marginBottom: 13,
  },
  button: {
    // textAlign: 'center',
    // marginBottom: 10,
    position: 'absolute', 
    left: 86,
    right: 86, //position, left, right together controls the button's length and horizontal location
    alignSelf: 'center', 
  },
  cardActions: {
    justifyContent: 'center', 
    alignItems: 'center',
    padding: 20,
  },
});

export default PostDetailsScreen;
